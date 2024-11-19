from tonutils.client import TonapiClient
from tonutils.wallet import WalletV5R1
from tonutils.utils import Address

# API key for accessing the Tonapi (obtainable from https://tonconsole.com)
API_KEY = "AFGNBPOEICSMP6AAAAAP4ILWBUS5H7RZF333GEP4KUY6OBUNHZCIPSVNJLB7P2RGR34KMEY"

# Set to True for test network, False for main network
IS_TESTNET = True

mnemonic_phrase = "hybrid dignity ecology ordinary balcony tide glory end recycle embrace comfort vacant junior slim reflect merge vacant style transfer valid measure potato pride dinner"

# Mnemonic phrase used to connect the wallet
MNEMONIC: list[str] = mnemonic_phrase.split(" ")

print(MNEMONIC, 'MNEMONIC')

# The address of the recipient
DESTINATION_ADDRESS = "0QC_NZK-QEP9oYEu-qgzjgN6kUajC34Sd0_6pLxxqjKXGHHz"

# Optional comment to include in the forward payload
COMMENT = "Hello from tonutils!"

# Amount to transfer in TON
AMOUNT = 0.1


async def main() -> None:
    client = TonapiClient(api_key=API_KEY, is_testnet=IS_TESTNET)
    wallet, public_key, private_key, mnemonic = WalletV5R1.from_mnemonic(
        client, MNEMONIC)

    print(wallet.address.to_str(), mnemonic)
    base64_address = wallet.address.to_str()

    base64_address_alt = Address(wallet.address).to_str(is_bounceable=False)

    print(base64_address_alt, 'base64_address_alt')

    tx_hash = await wallet.transfer(
        destination=DESTINATION_ADDRESS,
        amount=AMOUNT,
        body=COMMENT,
    )

    print(f"Successfully transferred {AMOUNT} TON!")
    print(f"Transaction hash: {tx_hash}")


if __name__ == "__main__":
    import asyncio

    asyncio.run(main())
