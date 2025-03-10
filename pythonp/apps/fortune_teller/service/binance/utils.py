import time

from binance import Client
from os import getenv
import logging
from decimal import Decimal
import json
import re


def generate_and_validate_symbol(symbol):
    symbol = symbol.upper()
    if not symbol.endswith("USDT"):
        symbol += "USDT"
    if not symbol.endswith("USDT"):
        raise ValueError(f"Symbol '{symbol}' is invalid. It must end with 'USDT'.")
    return symbol


TESTNET = getenv('TESTNET', 'True').lower() in ('true', '1', 't')

# 获取Binance API的key和secret
binance_api_key = getenv('BINANCE_API_KEY')
binance_api_secret = getenv('BINANCE_API_SECRET')

if not binance_api_key or not binance_api_secret:
    logging.error("Binance API key/secret is not set in the environment")
    raise ValueError("Binance API key/secret is not set in the environment")

# 初始化Binance客户端
binance_client = Client(binance_api_key, binance_api_secret, testnet=TESTNET)


def generate_symbol(currency):
    return f"{currency.upper()}"


def get_binance_client() -> Client:
    """
    获取Binance客户端。

    :return: Client, Binance客户端
    """
    return binance_client


def get_all_prices(length=30, filter_symbol=None):
    """
    获取所有交易对的价格，并根据需要过滤包含特定符号的交易对。

    :param length: int, 返回的交易对数量，默认是30
    :param filter_symbol: str, 需要过滤的符号，如 "USDT"，默认不过滤
    :return: dict, 符合条件的交易对的价格
    """
    try:
        prices = binance_client.get_all_tickers()
        if filter_symbol:
            prices = [p for p in prices if filter_symbol in p['symbol']]
        selected_prices = prices[:length]
        result = {p['symbol'].replace("USDT", ""): p['price'] for p in selected_prices}
        return result
    except Exception as e:
        logging.error(f"Error fetching prices: {e}")
        return {"error": "Error fetching prices."}


def get_common_currency_price(currency) -> str:
    """
    获取常见币种的价格。

    :param currency: CommonCurrency, 币种的枚举类型
    :return: str, 币种的价格
    """
    try:
        symbol = generate_and_validate_symbol(currency)
        ticker = binance_client.get_symbol_ticker(symbol=symbol)
        return ticker['price']
    except Exception as e:
        logging.error(f"Error fetching price for {currency.name}: {e}")
        return f"Error fetching price for {currency.name}."


def convert_currency(from_currency: str, to_currency: str,
                     amount: float) -> tuple[float, float, float] | str:
    """
    将指定数量的一个币种转换为另一个币种，并以常规小数格式返回结果。

    :param from_currency: str, 卖出币种的字符串（不区分大小写）
    :param to_currency: str, 买入币种的字符串（不区分大小写）
    :param amount: float, 卖出的币种数量
    :return: str, 以常规小数格式显示的买入的币种数量
    """
    try:
        # 将输入的字符串转为大写并转换为对应的枚举类型
        from_currency_enum =  generate_and_validate_symbol(from_currency)
        to_currency_enum =  generate_and_validate_symbol(to_currency)

        # 获取卖出币种的价格
        from_price = float(get_common_currency_price(from_currency_enum))

        # 获取买入币种的价格
        to_price = float(get_common_currency_price(to_currency_enum))

        print(from_price, to_price, "USDT")

        # 计算买入的币种数量
        converted_amount = (amount * from_price) / to_price

        # 返回格式化后的结果，保留小数点后8位
        result = float(f"{converted_amount:.8f}")

        return result, from_price, to_price
    except KeyError:
        logging.error(f"Invalid currency: {from_currency} or {to_currency}")
        return "Invalid currency provided."
    except Exception as e:
        logging.error(
            f"Error converting from {from_currency} to {to_currency}: {e}")
        return "Conversion error."


def create_order(symbol: str,
                 side: str,
                 quantity: float,
                 type="MARKET") -> str:
    """
    创建测试订单。

    :param type:  str, 订单类型（'MARKET'或'LIMIT'）
    :param symbol: str, 币种符号（如'BTCUSDT'）
    :param side: str, 买入或卖出（'BUY'或'SELL'）
    :param quantity: float, 订单数量
    :return: str, 订单结果
    """
    try:
        symbol = generate_and_validate_symbol(symbol)
        if generate_symbol(side) == 'BUY':
            order = binance_client.create_order(symbol=symbol,
                                                side="BUY",
                                                type=type,
                                                quantity=quantity)

        elif generate_symbol(side) == 'SELL':
            # create_test_order create_order
            order = binance_client.create_order(symbol=symbol,
                                                side="SELL",
                                                type=type,
                                                quantity=quantity)
        else:
            return "Invalid order side provided. Use 'BUY' or 'SELL'."

        logging.info(f"Order created: {order}")
        return order
    except Exception as e:
        logging.error(f"Error creating test order for {symbol}: {e}")
        raise ValueError(f"Error creating test order for {symbol}.")


def get_order_info(symbol):
    symbol = generate_and_validate_symbol(symbol)
    return binance_client.get_orderbook_ticker(symbol=symbol)


def get_order_status(symbol, order_id):
    symbol = generate_and_validate_symbol(symbol)
    order_details = binance_client.get_order(symbol=symbol, orderId=order_id)
    return order_details


def get_fees_via_symbol(symbol):
    symbol = generate_and_validate_symbol(symbol)
    return binance_client.get_trade_fee(symbol=symbol)


def calculate_fee(symbol="BTC", order_type='MARKET', amount=1, price=None):
    """
    计算指定订单类型的交易费用

    参数:
    client (object): 交易平台的 API 客户端
    symbol (str): 交易对，例如 'BNBBTC'
    order_type (str): 订单类型，默认为 'MARKET'
    amount (float): 交易的数量，默认为 1
    price (float): 价格，仅限于限价单（'LIMIT'），默认为 None

    返回:
    float: 交易费用
    """
    # 获取交易费率
    symbol = generate_and_validate_symbol(symbol)
    fees = binance_client.get_trade_fee(symbol=symbol)
    taker_fee_rate = float(fees['takerCommission']) / 100
    maker_fee_rate = float(fees['makerCommission']) / 100

    # 计算交易金额
    if order_type == 'MARKET':
        if price is None:
            # 获取当前市场价格（假设通过 API 获取）
            ticker = binance_client.get_ticker(symbol=symbol)
            price = float(ticker['lastPrice'])
    elif order_type == 'LIMIT':
        if price is None:
            raise ValueError("限价单必须提供价格")
    else:
        raise ValueError("未知的订单类型")

    trade_amount = price * amount

    # 根据订单类型计算费用
    if order_type == 'MARKET':
        fee = trade_amount * taker_fee_rate
    elif order_type == 'LIMIT':
        fee = trade_amount * maker_fee_rate

    return fee


def get_result_info(result):
    # 提取并检查字段
    symbol = result.get('symbol', 'UNKNOWN')
    order_id = result.get('orderId', -1)
    transact_time = result.get('transactTime', int(time.time() * 1000))
    cummulative_quote_qty = result.get('cummulativeQuoteQty', '0.0')
    status = result.get('status', 'UNKNOWN')
    order_type = result.get('type', 'UNKNOWN')
    side = result.get('side', 'UNKNOWN')

    # 数据类型转换和验证
    try:
        order_id = int(order_id)
    except ValueError:
        order_id = -1

    try:
        transact_time = int(transact_time)
    except ValueError:
        transact_time = int(time.time() * 1000)

    try:
        cummulative_quote_qty = float(cummulative_quote_qty)
    except ValueError:
        cummulative_quote_qty = 0.0

    # 提取并累加 fills 中的 commission 信息
    total_commission = 0.0
    commission_details = []
    fills = result.get('fills', [])
    for fill in fills:
        commission = fill.get('commission', '0.0')
        commission_asset = fill.get('commissionAsset', 'UNKNOWN')
        try:
            commission = float(commission)
        except ValueError:
            commission = 0.0
        total_commission += commission
        commission_details.append((commission, commission_asset))

    return symbol, order_id, status, order_type, side, transact_time, cummulative_quote_qty, total_commission, commission_details


def get_min_trade_quantity(symbol: str) -> Decimal:
    """
    获取指定交易对的最小交易量。

    :param symbol: str, 交易对符号（如'BTC'）
    :return: Decimal, 最小交易量
    """
    try:
        print(symbol, 'symbol')
        symbol = generate_and_validate_symbol(symbol)
        print(symbol, 'symbol')
        exchange_info = binance_client.get_exchange_info()
        for symbol_info in exchange_info['symbols']:
            if symbol_info['symbol'] == symbol:
                for filter in symbol_info['filters']:
                    if filter['filterType'] == 'LOT_SIZE':
                        return Decimal(filter['minQty']) * Decimal(10)
        return None
    except Exception as e:
        logging.error(f"Error fetching min trade quantity for {symbol}: {e}")
        return None


def is_min_trade_quantity_limit(symbol: str, quantity) -> bool:
    try:
        minimum = get_min_trade_quantity(symbol)
        return Decimal(quantity) < minimum, float(minimum)
    except  Exception as e:
        logging.error(f"Error checking min trade quantity limit for {symbol}: {e}")
        raise ValueError("This asset is not supported!")


def extract_json_from_markdown(markdown_text):
    """
    从Markdown文本中提取JSON数据。

    :param markdown_text: str, 包含JSON数据的Markdown文本
    :return: str, 提取出的JSON字符串
    """
    # 使用正则表达式匹配JSON数据
    json_match = re.search(r'```json(.*?)```', markdown_text, re.DOTALL)
    if json_match:
        return json_match.group(1).strip()
    else:
        return None


def extract_key_value_pairs(markdown_text):
    """
    从Markdown文本中提取单个键值对并构建JSON对象。

    :param markdown_text: str, 包含键值对的Markdown文本
    :return: dict, 构建的JSON对象
    """
    # 匹配单个键值对
    token_match = re.search(r'"token":\s*"([^"]+)"', markdown_text)
    action_match = re.search(r'"action":\s*"([^"]+)"', markdown_text)
    amount_match = re.search(r'"amount":\s*"([^"]+)"', markdown_text)

    data = {}
    if token_match:
        data["token"] = token_match.group(1)
    if action_match:
        data["action"] = action_match.group(1)
    if amount_match:
        data["amount"] = amount_match.group(1)

    return data if data else None


def process_recommendation(markdown_text):
    """
    从Markdown文本中提取JSON数据，去掉token中的USDT并输出token和action。

    :param markdown_text: str, 包含JSON数据的Markdown文本
    :return: tuple, 包含处理后的token、action和amount
    """
    # 提取JSON数据
    json_data = extract_json_from_markdown(markdown_text)

    if json_data:
        try:
            # 解析JSON数据
            data = json.loads(json_data)
        except json.JSONDecodeError:
            return "Error: Invalid JSON data"
    else:
        # 尝试提取单个键值对
        data = extract_key_value_pairs(markdown_text)
        if not data:
            return "Error: No JSON data found in the Markdown text"

    # 获取并处理token
    token = data.get("token", "")
    token = token.replace("USDT", "")

    # 获取action和amount
    action = data.get("action", "")
    amount = data.get("amount", 0)

    # 输出处理后的token、action和amount
    return token, action, amount


