from os import getenv
import logging

from pytonconnect import TonConnect

from pythonp.apps.tokenfate.service.ton.tc_storage import TcStorage

def get_connector(chat_id: int):
    MANIFEST_URL = getenv("MANIFEST_URL")

    print(f'================= get_connector ==============={MANIFEST_URL}')

    return TonConnect(MANIFEST_URL, storage=TcStorage(chat_id))
