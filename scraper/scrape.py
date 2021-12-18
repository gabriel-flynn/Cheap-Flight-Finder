import logging
import time
from concurrent.futures import ThreadPoolExecutor
import grpc

from flight_scraping_pb2 import SouthwestHeadersResponse
from flight_scraping_pb2_grpc import FlightScraperServicer, add_FlightScraperServicer_to_server
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from xvfbwrapper import Xvfb
import undetected_chromedriver_modified.v2 as uc

headers = {}
headers_last_fetched_at = None


def get_headers(event_data):
    if 'page/flights/products' in event_data['params']['request']['url']:
        global headers, headers_last_fetched_at
        headers = event_data['params']['request']['headers']
        headers_last_fetched_at = time.time()


def goToSouthwestWebsite():
    vdisplay = Xvfb(width=800, height=1280)
    vdisplay.start()

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
    try:
        element_present = EC.presence_of_element_located((By.CLASS_NAME, "form-field--placeholder"))
        WebDriverWait(driver, 60)

        from_btn = driver.find_element(By.CLASS_NAME, 'form-field--placeholder')
        from_btn.click()
        time.sleep(2)

        dal = driver.find_element(By.XPATH, "//span[contains(text(), '- DAL')]")
        dal.click()
        time.sleep(3)

        to_btn = driver.find_element(By.XPATH, "//div[text()='To']")
        to_btn.click()

        time.sleep(4)
        las = driver.find_element(By.XPATH, "//span[contains(text(), '- LAS')]")
        las.click()

        submit = driver.find_element(By.XPATH, "//button[@type='submit']")
        submit.click()
        time.sleep(10)

    except Exception as e:
        print(f"Encounter error: {e}")

    driver.close()
    vdisplay.stop()


class FlightScrapingServer(FlightScraperServicer):
    def GetSouthwestHeaders(self, request, context):
        # Only fetch new headers if it's been more than 5 minutes
        if headers_last_fetched_at is None or time.time() - headers_last_fetched_at > 259:
            goToSouthwestWebsite()
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
