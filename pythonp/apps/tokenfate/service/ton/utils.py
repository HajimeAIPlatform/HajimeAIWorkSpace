from os import getenv
from pytoniq_core import Address


async def calculate_fee(amount):
    try:
        fee_percentage = int(getenv('FEE_PERCENTAGE', 3))
        fee = amount * fee_percentage / 100

        # # 估算Gas费用
        # total_gas_fee_in_ton =  estimate_gas_fee()
        total_fee = fee
        return total_fee
    except Exception as e:
        print(f"Error estimating fee: {e}")
        return None


def estimate_gas_fee():
    try:
        FIXED_GAS_FEE = float(getenv('FIXED_GAS_FEE', 0.0055))
        return FIXED_GAS_FEE
    except Exception as e:
        print(f"Error estimating gas: {e}")
        return None, None, None


def get_bot_ton_address():
    ton_address = getenv('BOT_TON_ADDRESS')
    if not ton_address:
        raise ValueError("BOT_TON_ADDRESS environment variable not set")
    return ton_address


def convert_address_to_hex(address):
    try:
        return Address(address).to_str(is_bounceable=False)
    except Exception as e:
        print(f"Error converting address: {e}")
        return None
