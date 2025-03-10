# binance_module.py
import logging
from pythonp.apps.fortune_teller.service.binance.utils import (convert_currency,
                               get_binance_client, get_common_currency_price,
                               create_order, get_all_prices)

# 初始化Binance客户端
binance_client = get_binance_client()


def handle_binance_command(command):
    """
    处理与Binance相关的命令。

    :param command: str, 输入的命令
    :return: str, 命令的处理结果
    """
    logging.info(f"Handling command: {command}")
    if command == '/price':
        return get_all_prices()
