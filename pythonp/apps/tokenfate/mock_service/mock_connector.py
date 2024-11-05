
import asyncio
import random
import string
from pytoniq_core import Cell, TvmBitarray


def generate_valid_boc():
    # 假设 bits 是一个简单的字节序列
    bits = TvmBitarray(1023)
    bits.append(0)  # 添加一个简单的位作为示例
    refs = []

    # 创建 Cell 实例
    cell = Cell(bits, refs)

    # 生成 BOC 数据
    boc = cell.to_boc()
    return boc


async def send_transaction(*args,**kwargs):
    # 模拟一些延迟，类似于实际的网络请求
    await asyncio.sleep(random.uniform(0.5, 2.0))

    # 模拟生成一个随机的 BOC（Bag of Cells）字符串
    boc = generate_valid_boc()
    print(f"Generated BOC: {boc}")
    # 模拟返回的结果
    result = {
        "boc": boc,
        "status": "success",
        "transaction_id": ''.join(random.choices(string.ascii_letters + string.digits, k=16))
    }

    return result


