import asyncio
from pytoncenter import get_client
from pytoncenter.v3.models import GetAccountRequest

api_key = "be983dfbeee7922d319d8639af579d8a0898a4c1582dd4608574c92713850c49"
network = "testnet"

mnemonic_phrase = "hybrid dignity ecology ordinary balcony tide glory end recycle embrace comfort vacant junior slim reflect merge vacant style transfer valid measure potato pride dinner"

client = get_client(version="v3", network=network, api_key=api_key)
my_address = "0QAB-jMxbkKvOfqV7j0FrmzOJnazf6bp1WGTVWMcxcR8MGW3"
destination_address = "0QBkGnwEw17ir3c3-FlvSx5g4fjIYIPvZ0btUbuDXcw3Kw5N"


async def get_balance():
    try:
        account = await client.get_account(
            GetAccountRequest(address=my_address))
    except Exception as e:
        print("Error fetching account:", e)
        return

    print("=== Account Info ===")
    print("Symbol", "TON", "Balance:", account.balance / 1e9)


async def main():
    await get_balance()


if __name__ == "__main__":
    asyncio.run(main())
