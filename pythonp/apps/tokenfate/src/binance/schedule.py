from datetime import datetime, timedelta
import schedule
import threading
import random
import time
from binance import Client
import logging
import json
from src.binance.utils import get_binance_client
import os

def get_symbol_historical_prices(symbol, days):
    """
    获取特定交易对最近几天的价格数据。

    :param symbol: str, 交易对符号，如 "BTCUSDT"
    :param days: int, 获取数据的天数
    :return: list, 包含最近几天的收盘价
    """
    try:
        binance_client = get_binance_client()
        # 定义时间间隔
        interval = Client.KLINE_INTERVAL_1DAY

        # 获取当前时间和指定天数前的时间
        end_time = datetime.now()
        start_time = end_time - timedelta(days=days)

        logging.info(f"Fetching data for {symbol} from {start_time} to {end_time}")

        end_timestamp = int(end_time.timestamp() * 1000)
        start_timestamp = int(start_time.timestamp() * 1000)

        # 获取 K 线数据
        klines = binance_client.get_historical_klines(symbol, interval, start_timestamp,
                                                      end_timestamp)

        # 只获取收盘价
        close_prices = [float(kline[4]) for kline in klines]

        logging.info(f"Fetched {len(close_prices)} prices for {symbol}")

        return close_prices
    except Exception as e:
        logging.error(f"Error fetching historical prices for {symbol}: {e}")
        return []


def fetch_and_store_usdt_historical_prices(days=7):
    """
    获取所有USDT交易对最近几天的价格数据，并存储到本地文件。

    :param days: int, 获取数据的天数，默认是7天
    """
    try:
        # 获取所有交易对的价格
        binance_client = get_binance_client()
        prices = binance_client.get_all_tickers()
        usdt_symbols = [p['symbol'] for p in prices if p['symbol'].endswith('USDT')]

        result = {}
        for symbol in usdt_symbols:
            result[symbol] = get_symbol_historical_prices(symbol, days)

        # 存储到本地文件
        with open('usdt_historical_prices.json', 'w') as f:
            json.dump(result, f)

        logging.info("Historical prices fetched and stored successfully.")
    except Exception as e:
        logging.error(f"Error fetching historical prices for USDT pairs: {e}")


def get_random_usdt_historical_prices(sample_size=30):
    """
    从本地文件中读取USDT交易对的价格数据，并随机选择30个symbol返回。

    :param sample_size: int, 返回的交易对数量，默认是30个
    :return: dict, 包含随机选择的USDT交易对的价格数据
    """
    try:
        # 检查文件是否存在
        if not os.path.exists('usdt_historical_prices.json'):
            logging.error("Historical prices file does not exist.")
            return {"error": "Historical prices file does not exist."}

        # 读取本地文件
        with open('usdt_historical_prices.json', 'r') as f:
            data = json.load(f)

        # 随机选择sample_size个symbol
        selected_symbols = random.sample(list(data.keys()), sample_size)
        result = {symbol: data[symbol] for symbol in selected_symbols}

        return result
    except Exception as e:
        logging.error(f"Error reading historical prices from file: {e}")
        return {"error": "Error reading historical prices from file."}


def run_schedule():
    """
    运行定时任务的线程函数。
    """
    while True:
        schedule.run_pending()
        time.sleep(1)


def start_schedule_thread():
    """
    启动定时任务线程。
    """
    # 设置定时任务，每天0点执行
    schedule.every().day.at("00:00").do(fetch_and_store_usdt_historical_prices)
    # schedule.every(1).minutes.do(fetch_and_store_usdt_historical_prices)
    # 启动定时任务线程
    schedule_thread = threading.Thread(target=run_schedule)
    schedule_thread.daemon = True
    schedule_thread.start()
