# export PYTHONPATH=$(pwd)
# $env:PYTHONPATH = (Get-Location).Path
from dotenv import load_dotenv

load_dotenv()

from src.binance.utils import (get_binance_client, generate_and_validate_symbol,
                               get_common_currency_price, create_order,
                               get_order_info, get_all_prices, get_min_trade_quantity, is_min_trade_quantity_limit,
                               process_recommendation)

from src.binance.schedule import get_random_usdt_historical_prices,fetch_and_store_usdt_historical_prices

def test(user_input):
    try:
        # 将用户输入的币种转换为对应的枚举
        currency = generate_and_validate_symbol(user_input)

        # 获取该币种的价格
        price = get_common_currency_price(currency)
        print(f"The current price of {user_input} is {price} USDT.")
    except KeyError:
        print(
            f"Invalid currency: {user_input}. Please enter a valid cryptocurrency symbol like BTC, ETH, or BNB."
        )


if __name__ == "__main__":
    fetch_and_store_usdt_historical_prices()
    result = get_random_usdt_historical_prices()
    print(result)
    # result = get_all_prices(10000,"USDT")
    # print(result,len(result),'get_all_prices')

    # usdt_historical_prices = get_all_usdt_historical_prices(3)
    # print(usdt_historical_prices)

    # markdown_text = """
    # 今天推荐的交易建议如下：
    # {
    #   "token": "ADAUSDT",
    #   "action": "买入"
    # }
    # ```
    # """
    # result = process_recommendation(markdown_text)
    # if isinstance(result, tuple):
    #     token, action, ammount = result
    #     print(f"Token: {token}")
    #     print(f"Action: {action}")
    # else:
    #     print(result)

    # test("BTC")
    # test("TON")
    # test("ETH")
    # test("BNB")
    #
    # # 示例用法: 买0.5个BTC，用TON支付
    # amount_btc = 1
    # converted_ton = convert_currency("BTC", "TON", amount_btc)
    # print(f"To buy {amount_btc} BTC, you need {converted_ton} TON.")
    #
    # amount_ton = float(converted_ton)
    # converted_ton = convert_currency("TON", "BTC", amount_ton)
    # print(f"To sell {converted_ton} TON, you need {converted_ton} BTC.")

    # res = get_min_trade_quantity("DOGE")
    # print(res)

    # buy_result =  is_min_trade_quantity_limit("BTC", 0.00001)
    # print(buy_result)
    #
    # # 买入示例
    # buy_result = create_order("BTC", "BUY", float(res))  # 假设买入0.001 BTC
    # print(buy_result)
    #
    #
    # # 卖出示例
    # sell_result = create_order("BTC", "SELL", 0.00001)  # 假设卖出0.001 BTC
    # print(sell_result)
