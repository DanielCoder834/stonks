# This is a sample Python script.
import csv
import threading
from typing import List

import redis
# Press ⌃R to execute it or replace it with your code.
# Press Double ⇧ to search everywhere for classes, files, tool windows, actions, and settings.
import requests

from alpha_vantage_api import AlphaVantageAPI
from resp_client import start_resp

r = redis.StrictRedis(host='localhost', port=6379, decode_responses=True)


# Returns an array of array list in this case a CSV file
def decode_csv_to_2d_arr(csv_str):
    with requests.Session() as s:
        download = s.get(csv_str)
        decoded_content = download.content.decode('utf-8')
        cr = csv.reader(decoded_content.splitlines(), delimiter=',')
        # 2-D Array List
        csv_list = list(cr)
        return csv_list


def connect_to_db():
    r.execute_command('set others test2')
    data = r.execute_command('get others')
    print(data)


if __name__ == '__main__':
    # t1 = threading.Thread(connect_to_db())
    # t2 = threading.Thread(connect_to_db())
    # t1.start()
    # t2.start()
    # AlphaVantageAPI()
    start_resp(r)
