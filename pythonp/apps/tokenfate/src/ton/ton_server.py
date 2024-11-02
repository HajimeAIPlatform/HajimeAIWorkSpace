from tonutils.client import TonapiClient
from tonutils.wallet import WalletV5R1, WalletV4R2, WalletV4R1, WalletV3R1, WalletV3R2
from tonutils.utils import Address
from os import getenv
import logging

# API key for accessing the Tonapi (obtainable from https://tonconsole.com)
API_KEY = getenv('TON_API_KEY', '')
TESTNET = getenv('TESTNET', 'True').lower() in ('true', '1', 't')
mnemonic_phrase = getenv('TON_MNEMONIC', '')
# Mnemonic phrase used to connect the wallet
MNEMONIC: list[str] = mnemonic_phrase.split(" ")

# The address of the recipient
# DESTINATION_ADDRESS = "0QBkGnwEw17ir3c3-FlvSx5g4fjIYIPvZ0btUbuDXcw3Kw5N"

# Optional comment to include in the forward payload
# COMMENT = "Hello from tonutils!"

# Amount to transfer in TON
# AMOUNT = 0.01


def get_client():
    return TonapiClient(api_key=API_KEY, is_testnet=TESTNET)


async def get_ton_wallet_address() -> Address:
    client = get_client()
    wallet, public_key, private_key, mnemonic = WalletV4R2.from_mnemonic(
        client, MNEMONIC)
    return Address(wallet.address).to_str(is_bounceable=False)


async def send_transaction(destination_address: str, amount: float,
                           comment: str) -> None:
    client = get_client()
    wallet, public_key, private_key, mnemonic = WalletV5R1.from_mnemonic(
        client, MNEMONIC)
    # Build the transaction
    tx_hash = await wallet.transfer(destination=destination_address,
                                    amount=amount,
                                    body=comment)
    logging.info(
        f"send_transaction: destination_address: {destination_address}, amount: {amount}, comment: {comment} Ton, Transaction hash: {tx_hash}"
    )
    return tx_hash
