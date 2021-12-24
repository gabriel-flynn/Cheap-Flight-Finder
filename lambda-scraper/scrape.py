import json
import re

from selenium import webdriver

chrome = None


def handler(event, context=None):
    global chrome

    if chrome is None:
        options = webdriver.ChromeOptions()
        options.binary_location = '/opt/chrome/chrome'
        options.add_argument('--headless')
        options.add_argument('--no-sandbox')
        options.add_argument("--disable-gpu")
        options.add_argument("--window-size=1280x1696")
        options.add_argument("--single-process")
        options.add_argument("--disable-dev-shm-usage")
        options.add_argument("--disable-dev-tools")
        options.add_argument("--no-zygote")
        options.add_argument("--user-data-dir=/tmp/chrome-user-data")
        options.add_argument("--remote-debugging-port=9222")
        chrome = webdriver.Chrome("/opt/chromedriver",
                                  options=options)

    chrome.get(event['url'])
    json_str = re.search("FlightData = '({.*})';", chrome.page_source).group(1)
    return json.loads(json_str.replace("&quot;", "\""))
