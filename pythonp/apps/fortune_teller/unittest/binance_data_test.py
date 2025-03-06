get_order_status = {
    "symbol":
    "BNBBTC",  # 交易对
    "orderId":
    28,  # 订单ID
    "orderListId":
    -1,  # 订单列表ID（如果不是OCO订单，通常为-1）
    "clientOrderId":
    "6gCrw2kRUAF9CvJDGP16IP",  # 客户端订单ID
    "price":
    "0.00000000",  # 订单价格（市价订单通常为0）
    "origQty":
    "100.00000000",  # 原始订单数量
    "executedQty":
    "100.00000000",  # 已执行数量
    "cummulativeQuoteQty":
    "0.00299100",  # 累计报价数量
    "status":
    "FILLED",  # 订单状态，例如NEW, PARTIALLY_FILLED, FILLED, CANCELED
    "timeInForce":
    "GTC",  # 时间有效性
    "type":
    "MARKET",  # 订单类型
    "side":
    "BUY",  # 订单方向
    "stopPrice":
    "0.00000000",  # 止损价格（如果适用）
    "icebergQty":
    "0.00000000",  # 冰山订单数量（如果适用）
    "time":
    1507725176595,  # 订单创建时间（时间戳）
    "updateTime":
    1507725176595,  # 订单更新时间（时间戳）
    "isWorking":
    True,  # 订单是否正在工作
    "origQuoteOrderQty":
    "0.00000000",  # 原始报价订单数量
    "fills": [  # 填充信息（交易详情）
        {
            "price": "0.00002991",  # 成交价格
            "qty": "100.00000000",  # 成交数量
            "commission": "0.00000000",  # 佣金
            "commissionAsset": "BTC",  # 佣金资产
            "tradeId": 34  # 交易ID
        }
    ]
}

get_fees_via_symbol = {
    "tradeFee": [{
        "symbol": "BNBBTC",  # 交易对
        "maker": 0.001,  # Maker 费率
        "taker": 0.001  # Taker 费率
    }],
    "success":
    True  # 请求是否成功
}

# {'symbol': 'BTCUSDT', 'orderId': 15534692, 'orderListId': -1, 'clientOrderId': 'm3yR023W0EX90NlnRgkLdX', 'transactTime': 1725436482908, 'price': '0.0000
# 0000', 'origQty': '0.00100000', 'executedQty': '0.00100000', 'cummulativeQuoteQty': '56.81984000', 'status': 'FILLED', 'timeInForce': 'GTC', 'type': 'MA
# RKET', 'side': 'BUY', 'workingTime': 1725436482908, 'fills': [{'price': '56819.84000000', 'qty': '0.00100000', 'commission': '0.00000000', 'commissionAsset': 'BTC', 'tradeId': 3323905}], 'selfTradePreventionMode': 'EXPIRE_MAKER'}
# {'symbol': 'BTCUSDT', 'orderId': 15534697, 'orderListId': -1, 'clientOrderId': 'JQU0Dj0AlwRXdnEVyex1UD', 'transactTime': 1725436483391, 'price': '0.0000
# 0000', 'origQty': '0.00100000', 'executedQty': '0.00100000', 'cummulativeQuoteQty': '56.81952000', 'status': 'FILLED', 'timeInForce': 'GTC', 'type': 'MA
# RKET', 'side': 'SELL', 'workingTime': 1725436483391, 'fills': [{'price': '56819.52000000', 'qty': '0.00100000', 'commission': '0.00000000', 'commissionAsset': 'USDT', 'tradeId': 3323906}], 'selfTradePreventionMode': 'EXPIRE_MAKER'}
