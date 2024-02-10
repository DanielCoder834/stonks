from typing import List

import requests
from dotenv import load_dotenv

load_dotenv()

def bool_cast_pure_string(boolean: bool):
    if boolean:
        return "true"
    else:
        return "false"


def create_url(dictionary_of_values: dict):
    url = ""
    for key, value in dictionary_of_values:
        if key == "basic_start":
            url += value
        elif value is None:
            url += ""
        elif type(value) == List:
            string_of_var_list = '&' + key + "=" + ','.join(value)
            url += string_of_var_list
        elif type(value) is bool:
            string_of_var_bool = "&" + key + "=" + bool_cast_pure_string(value)
            url += string_of_var_bool
        else:
            if not isinstance(value, str):
                raise ValueError("Value must be a string, something has gone wrong")
            string_of_var = "&" + key + "=" + value
            url += string_of_var
    return url


class AlphaVantageAPI:
    def __init__(self):
        self.vantage_api_key = os.getenv('ALPHA_VANTAGE_KEY')

    # Takes in an url and returns the json object of that result
    def request_new_data(self, url):
        request_url = url + self.vantage_api_key
        r = requests.get(request_url)
        return r.json()

    def intraday(self, symbols: List[str], interval: str, adjusted: bool = True, extended_hours: bool = True,
                 month: str = None, outputsize: str = None, datasize: str = None):
        if 5 < len(symbols):
            raise TypeError("I don't have premium so only a max of 5")
        return create_url({"basic_start": "https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY",
                           "SYMBOLS": symbols, "INTERVAL": interval, "adjusted": adjusted,
                           "extended_hours": extended_hours,
                           "month": month, "outputsize": outputsize, "datasize": datasize,
                           "apikey": self.vantage_api_key})

    # Returns an url for the news
    def news_segment(self, tickers: str = None, topics: str = None, time_from: str = None,
                     time_to: str = None, sort: str = None,
                     limit: str = None):
        basic_url_start = "https://www.alphavantage.co/query?"
        function = "NEWS_SENTIMENT"
        if type(time_to) != type(time_from):
            raise ValueError("Time to and from must have same types")
        tickers_str = "" if (tickers is None) else '&tickers=' + tickers
        topics_str = "" if (topics is None) else '&topics=' + topics
        time_from_str = "" if (time_from is None) else '&time_from=' + time_from
        time_to_str = "" if (time_to is None) else '&time_to=' + time_to
        sort_str = "" if (sort is None) else '&sort=' + sort
        limit_str = "" if (limit is None) else '&limit=' + limit
        api_str = '&apikey=' + self.vantage_api_key
        url_str = basic_url_start + 'function=' + function + tickers_str + topics_str + time_from_str + \
                  time_to_str + sort_str + limit_str + api_str
        return url_str

    # Returns an url for top 20 gainers or top 20 losers
    def losers_gainers(self):
        function = "TOP_GAINERS_LOSERS"
        basic_url_start = 'https://www.alphavantage.co/query?'
        api_str = '&apikey=' + self.vantage_api_key
        return basic_url_start + 'function=' + function + api_str

    # Abstract out the code in both analytic window methods
    def analytics_fixed_window(self, symbols: List[str], date_range: List[str], interval: str, calculation: List[str],
                               ohlc: str = None):
        if 1 > len(date_range) > 3:
            raise TypeError("date_range must have a start and end time for the range")
        if 5 < len(symbols):
            raise TypeError("I don't have premium so only a max of 5")
        basic_url_start = 'https://alphavantageapi.co/timeseries/running_analytics?'
        symbols_str = '&SYMBOLS=' + ','.join(symbols)
        date_range_str = '&RANGE='.join(date_range)
        interval_str = '&INTERVAL=' + interval
        calculation_str = '&CALCULATIONS=' + ','.join(calculation)
        ohlc_str = "" if (ohlc is None) else '&OHLC=' + ohlc
        api_str = '&apikey=' + self.vantage_api_key
        return basic_url_start + symbols_str + date_range_str + interval_str + ohlc_str + calculation_str + api_str

    def analytics_sliding_window(self, symbols: List[str], date_range: List[str], interval: str,
                                 window_size: int, calculation: List[str], ohlc: str = None):
        # If it is less than 1 or greater than 3
        if 1 > len(date_range) > 3:
            raise TypeError("date_range must have a start and end time for the range")
        if 5 < len(symbols):
            raise TypeError("I don't have premium so only a max of 5")
        basic_url_start = 'https://alphavantageapi.co/timeseries/running_analytics?'
        symbols_str = '&SYMBOLS=' + ','.join(symbols)
        date_range_str = '&RANGE='.join(date_range)
        interval_str = '&INTERVAL=' + interval
        window_size_str = '&WINDOW_SIZE=' + str(window_size)
        calculation_str = '&CALCULATIONS=' + ','.join(calculation)
        ohlc_str = "" if (ohlc is None) else '&OHLC=' + ohlc
        api_str = '&apikey=' + self.vantage_api_key
        url_str = basic_url_start + symbols_str + date_range_str + interval_str + \
                  ohlc_str + window_size_str + calculation_str + api_str
        return url_str

    def company_overview(self, symbols: List[str]):
        basic_url_start = 'https://www.alphavantage.co/query?function=OVERVIEW'
        symbols_str = '&SYMBOLS=' + ','.join(symbols)
        api_str = '&apikey=' + self.vantage_api_key
        return basic_url_start + symbols_str + api_str

    def income_statement(self, symbols: List[str]):
        basic_url_start = 'https://www.alphavantage.co/query?function=INCOME_STATEMENT'
        symbols_str = '&SYMBOLS=' + ','.join(symbols)
        api_str = '&apikey=' + self.vantage_api_key
        return basic_url_start + symbols_str + api_str

    def balance_sheet(self, symbols: List[str]):
        basic_url_start = 'https://www.alphavantage.co/query?function=BALANCE_SHEET'
        symbols_str = '&SYMBOLS=' + ','.join(symbols)
        api_str = '&apikey=' + self.vantage_api_key
        return basic_url_start + symbols_str + api_str

    def cash_flow(self, symbols: List[str]):
        basic_url_start = 'https://www.alphavantage.co/query?function=CASH_FLOW'
        symbols_str = '&SYMBOLS=' + ','.join(symbols)
        api_str = '&apikey=' + self.vantage_api_key
        return basic_url_start + symbols_str + api_str

    def earnings(self, symbols: List[str]):
        basic_url_start = 'https://www.alphavantage.co/query?function=EARNINGS'
        symbols_str = '&SYMBOLS=' + ','.join(symbols)
        api_str = '&apikey=' + self.vantage_api_key
        return basic_url_start + symbols_str + api_str

    # Returns the url for the listed companies
    def listing_status(self, date: str = None, state: str = None):
        basic_url_start = 'https://www.alphavantage.co/query?function=LISTING_STATUS'
        date_str = "" if (date is None) else '&date=' + date
        state_str = "" if (state is None) else '&state=' + state
        api_str = '&apikey=' + self.vantage_api_key
        csv_str = basic_url_start + date_str + state_str + api_str
        return csv_str

    # Returns the url for the earning calender of a company
    def earning_calender(self, symbol: str = None, horizon: str = None):
        basic_url_start = "https://www.alphavantage.co/query?function=EARNINGS_CALENDAR"
        symbol_str = "" if (symbol is None) else '&date=' + symbol
        horizon_str = "" if (horizon is None) else '&state=' + horizon
        api_str = '&apikey=' + self.vantage_api_key
        csv_str = basic_url_start + symbol_str + horizon_str + api_str
        return csv_str

    def ipo_calender(self):
        basic_url_start = "https://www.alphavantage.co/query?function=IPO_CALENDAR"
        api_str = '&apikey=' + self.vantage_api_key
        return basic_url_start + api_str
