import logging
import threading
import time
from concurrent.futures import ThreadPoolExecutor
import grpc
from selenium import webdriver

from flight_scraping_pb2 import SouthwestHeadersResponse
from flight_scraping_pb2_grpc import FlightScraperServicer, add_FlightScraperServicer_to_server
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from xvfbwrapper import Xvfb
import undetected_chromedriver_modified.v2 as uc

headers = {}
headers_last_fetched_at = None
sem = threading.Semaphore()


def get_headers(event_data):
    if 'page/flights/products' in event_data['params']['request']['url']:
        global headers, headers_last_fetched_at
        headers = event_data['params']['request']['headers']
        headers_last_fetched_at = time.time()


def goToSouthwestWebsite():
    # vdisplay = Xvfb(width=800, height=1280)
    # vdisplay.start()

    options = uc.ChromeOptions()
    options.add_argument(f'--no-first-run --no-service-autorun --password-store=basic')
    options.user_data_dir = f'./tmp/test_undetected_chromedriver'
    options.add_argument(f'--disable-gpu')
    options.add_argument(f'--no-sandbox')
    options.add_argument(f'--disable-dev-shm-usage')
    driver = uc.Chrome(
        options=options,
        headless=False,
        enable_cdp_events=True)
    driver.add_cdp_listener('Network.requestWillBeSent', get_headers)
    driver.get('https://mobile.southwest.com/air/booking/shopping')
    time.sleep(15)
    try:
        element_present = EC.presence_of_element_located((By.CLASS_NAME, "form-field--placeholder"))
        WebDriverWait(driver, 60)
        driver.execute_cdp_cmd('Network.setUserAgentOverride', {"userAgent":"Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.104 Mobile Safari/537.36"})
        action = webdriver.ActionChains(driver)

        from_btn = driver.find_element(By.CLASS_NAME, 'form-field--placeholder')
        action.move_to_element(from_btn)
        action.click().perform()
        time.sleep(2)

        dal = driver.find_element(By.XPATH, "//span[contains(text(), '- ALB')]")
        action.move_to_element(dal)
        action.click().perform()
        time.sleep(3)

        to_btn = driver.find_element(By.XPATH, "//div[text()='To']")
        action.move_to_element(to_btn)
        action.click().perform()

        time.sleep(4)
        las = driver.find_element(By.XPATH, "//span[contains(text(), '- ABQ')]")
        action.move_to_element(las)
        action.click().perform()

        submit = driver.find_element(By.XPATH, "//button[@type='submit']")
        action.move_to_element(submit)
        time.sleep(1)
        action.click().perform()
        time.sleep(10)

    except Exception as e:
        print(f"Encounter error: {e}")

    driver.close()
    # vdisplay.stop()


class FlightScrapingServer(FlightScraperServicer):
    def GetSouthwestHeaders(self, request, context):
        # Only fetch new headers if it's been more than 5 minutes
        # if headers_last_fetched_at is None or time.time() - headers_last_fetched_at > 259:
        sem.acquire()
        goToSouthwestWebsite()
        sem.release()
        resp = SouthwestHeadersResponse(headers=headers)
        return resp


if __name__ == '__main__':
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(levelname)s - %(message)s',
    )
    logging.getLogger("urllib3").setLevel(logging.ERROR)

    server = grpc.server(ThreadPoolExecutor())
    add_FlightScraperServicer_to_server(FlightScrapingServer(), server)
    port = 9999
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    logging.info('server ready on port %r', port)
    server.wait_for_termination()
