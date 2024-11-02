import asyncio
import random
import string


async def send_transaction(*args,**kwargs):
    # 模拟一些延迟，类似于实际的网络请求
    await asyncio.sleep(random.uniform(0.5, 2.0))

    # 模拟生成一个随机的交易哈希
    tx_hash = "https://toncenter.com/api/v3/transactionsByMessage?direction=out&msg_hash=90aec8965afabb16ebc3cb9b408ebae71b618d78788bc80d09843593cac98da4&limit=128&offset=0"

    return tx_hash



