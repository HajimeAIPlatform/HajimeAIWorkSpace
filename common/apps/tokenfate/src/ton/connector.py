from pytonconnect import TonConnect
from src.ton.tc_storage import TcStorage
from os import getenv
import logging

def get_connector(chat_id: int):
    MANIFEST_URL = getenv("MANIFEST_URL")

    print(f'================= get_connector ==============={MANIFEST_URL}')

    return TonConnect(MANIFEST_URL, storage=TcStorage(chat_id))