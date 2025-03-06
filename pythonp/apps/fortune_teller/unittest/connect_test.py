from dotenv import load_dotenv

load_dotenv()


import asyncio
import json
import os
import logging
import time

from pytonconnect import TonConnect
from pytonconnect.storage import IStorage
import redis.asyncio as redis

from pythonp.apps.fortune_teller.service.ton.messages import get_comment_message



# 配置 Redis 客户端
REDIS_HOST = os.getenv("REDIS_HOST", "10.10.101.126")
REDIS_PORT = int(os.getenv("REDIS_PORT", 6389))
if not REDIS_HOST or not REDIS_PORT:
    raise Exception("REDIS_HOST and REDIS_PORT must be set")

client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT)

logging.info(f"Connected to Redis at {REDIS_HOST}:{REDIS_PORT}")


# 定义 TcStorage 类
class TcStorage(IStorage):

    def __init__(self, chat_id: int):
        self.chat_id = chat_id

    def _get_key(self, key: str):
        return f"{self.chat_id}:{key}"

    async def set_item(self, key: str, value: str):
        print(f"Setting key {key} to {value}")
        await client.set(name=self._get_key(key), value=value)

    async def get_item(self, key: str, default_value: str = None):
        value = await client.get(name=self._get_key(key))
        print(f"Getting key {key}")
        return value.decode() if value else default_value

    async def remove_item(self, key: str):
        print(f"Removing key {key}")
        await client.delete(self._get_key(key))


# 获取 TonConnect 实例
def get_connector(chat_id: int):
    MANIFEST_URL = os.getenv("MANIFEST_URL")

    print(f'================= get_connector ==============={MANIFEST_URL}')

    return TonConnect(MANIFEST_URL, storage=TcStorage(chat_id))


# 主函数
async def main():
    chat_id = 6969837043  # 你可以先写死这个值
    connector = get_connector(chat_id)
    connected = await connector.restore_connection()
    if not connected:
        return 'Connect wallet first!'

    print(connector.account.address)

    transaction = {
        'valid_until':
            int(time.time() + 3600),
        'messages': [
            get_comment_message(
                destination_address="0QB_A1-mv0OPJu6QbGClima_joZr_xlAfPcvE_E8pjH9XayK",
                amount=0.001 * 10 ** 9,
                comment="7777")
        ]
    }

    result = await asyncio.wait_for(
        connector.send_transaction(transaction=transaction), 300)
    print(result,'result')

async def save_session_data(chat_id, session_data,last_event_id=1725201255997214):
    storage = TcStorage(chat_id)
    await storage.set_item("last_event_id", last_event_id)
    # await storage.set_item("connection", json.dumps(session_data))
    await storage.set_item("connection", session_data)

# 示例：首次连接后保存会话数据
chat_id = 6969837043
session_data = {
  "type": "http",
  "connectEvent": {
    "id": 14,
    "event": "connect",
    "payload": {
      "items": [
        {
          "name": "ton_addr",
          "address": "0:bf3592be4043fda1812efaa8338e037a9146a30b7e12774ffaa4bc71aa329718",
          "network": "-3",
          "publicKey": "fb59a2cf01989beac4b1c63de1b00817acfa0915d76136257f4c5a8c548679ad",
          "walletStateInit": "te6cckECFgEAAwQAAgE0AgEAUQAAAAApqaMX+1mizwGYm+rEscY94bAIF6z6CRXXYTYlf0xajFSGea1AART/APSkE/S88sgLAwIBIAkEBPjygwjXGCDTH9Mf0x8C+CO78mTtRNDTH9Mf0//0BNFRQ7ryoVFRuvKiBfkBVBBk+RDyo/gAJKTIyx9SQMsfUjDL/1IQ9ADJ7VT4DwHTByHAAJ9sUZMg10qW0wfUAvsA6DDgIcAB4wAhwALjAAHAA5Ew4w0DpMjLHxLLH8v/CAcGBQAK9ADJ7VQAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AABwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbtIH+gDU1CL5AAXIygcVy//J0Hd0gBjIywXLAiLPFlAF+gIUy2sSzMzJc/sAyEAUgQEI9FHypwICAUgTCgIBIAwLAFm9JCtvaiaECAoGuQ+gIYRw1AgIR6STfSmRDOaQPp/5g3gSgBt4EBSJhxWfMYQCASAODQARuMl+1E0NcLH4AgFYEg8CASAREAAZrx32omhAEGuQ64WPwAAZrc52omhAIGuQ64X/wAA9sp37UTQgQFA1yH0BDACyMoHy//J0AGBAQj0Cm+hMYALm0AHQ0wMhcbCSXwTgItdJwSCSXwTgAtMfIYIQcGx1Z70ighBkc3RyvbCSXwXgA/pAMCD6RAHIygfL/8nQ7UTQgQFA1yH0BDBcgQEI9ApvoTGzkl8H4AXTP8glghBwbHVnupI4MOMNA4IQZHN0crqSXwbjDRUUAIpQBIEBCPRZMO1E0IEBQNcgyAHPFvQAye1UAXKwjiOCEGRzdHKDHrFwgBhQBcsFUAPPFiP6AhPLassfyz/JgED7AJJfA+IAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABmNjqp0="
        }
      ],
      "device": {
        "platform": "android",
        "appName": "Tonkeeper",
        "appVersion": "4.9.0",
        "maxProtocolVersion": 2,
        "features": [
          "SendTransaction",
          {
            "name": "SendTransaction",
            "maxMessages": 4
          }
        ]
      }
    }
  },
  "session": {
    "sessionKeyPair": {
      "publicKey": "22b726459a586aea14ca8cf4bd4fc998cd507b349fddf974d2dd64db97162408",
      "secretKey": "5d8be2ed0bd1fbe4a4fe969d74fb26a1f419aca950497fd0013b05fbee30a25b"
    },
    "walletPublicKey": "f3504259930a40d7a9fbc0682cc8413f87743c32f2bfbb2211d366220b907a2a",
    "bridgeUrl": "https://bridge.tonapi.io/bridge"
  },
  "lastWalletEventId": 14,
  "nextRpcRequestId": 0
}


session_data1 = {
    "type": "http",
    "session": {
        "session_private_key": "60425192262193d35e9819588cc4c3d3264ccff9ea68b99bdc3f46589665f843",
        "wallet_public_key": "0f41d645f06b7d9525499a3a605f51596a0e406f831dbc78184c12ed9d4fe04c",
        "bridge_url": "https://bridge.ton.space/bridge"
    },
    "last_wallet_event_id": 7,
    "connect_event": {
        "id": 7,
        "event": "connect",
        "payload": {
            "items": [
                {
                    "name": "ton_addr",
                    "address": "0:57322da2ee2b5a03dca3db388b5d07416c8d9e503d285558eccec8399e8f833d",
                    "network": "-239",
                    "publicKey": "964ffc82e10ccc4cd63f10f132feb18a10c66847b3e15fdbec0f3aadf747ce17",
                    "walletStateInit": "te6cckECFgEAArEAAgE0ARUBFP8A9KQT9LzyyAsCAgEgAw4CAUgEBQLc0CDXScEgkVuPYyDXCx8gghBleHRuvSGCEHNpbnS9sJJfA+CCEGV4dG66jrSAINchAdB01yH6QDD6RPgo+kQwWL2RW+DtRNCBAUHXIfQFgwf0Dm+hMZEw4YBA1yFwf9s84DEg10mBAoC5kTDgcOIREAIBIAYNAgEgBwoCAW4ICQAZrc52omhAIOuQ64X/wAAZrx32omhAEOuQ64WPwAIBSAsMABezJftRNBx1yHXCx+AAEbJi+1E0NcKAIAAZvl8PaiaECAoOuQ+gLAEC8g8BHiDXCx+CEHNpZ2668uCKfxAB5o7w7aLt+yGDCNciAoMI1yMggCDXIdMf0x/TH+1E0NIA0x8g0x/T/9cKAAr5AUDM+RCaKJRfCtsx4fLAh98Cs1AHsPLQhFEluvLghVA2uvLghvgju/LQiCKS+ADeAaR/yMoAyx8BzxbJ7VQgkvgP3nDbPNgRA/btou37AvQEIW6SbCGOTAIh1zkwcJQhxwCzji0B1yggdh5DbCDXScAI8uCTINdKwALy4JMg1x0GxxLCAFIwsPLQiddM1zkwAaTobBKEB7vy4JPXSsAA8uCT7VXi0gABwACRW+Dr1ywIFCCRcJYB1ywIHBLiUhCx4w8g10oSExQAlgH6QAH6RPgo+kQwWLry4JHtRNCBAUHXGPQFBJ1/yMoAQASDB/RT8uCLjhQDgwf0W/LgjCLXCgAhbgGzsPLQkOLIUAPPFhL0AMntVAByMNcsCCSOLSHy4JLSAO1E0NIAURO68tCPVFAwkTGcAYEBQNch1woA8uCO4sjKAFjPFsntVJPywI3iABCTW9sx4ddM0ABRgAAAAD///4jLJ/5BcIZmJmsfiHiZf1jFCGM0I9nwr+32B51W+6PnC6CODfLG"
                }
            ],
            "device": {
                "platform": "browser",
                "appName": "telegram-wallet",
                "appVersion": "1",
                "maxProtocolVersion": 2,
                "features": [
                    "SendTransaction",
                    {
                        "name": "SendTransaction",
                        "maxMessages": 4
                    }
                ]
            }
        }
    },
    "next_rpc_request_id": 0
}

def get_data(data):
    session_data1 = {
        "type": data.get("type", "http"),
        "session": {
            "session_private_key": data["session"]["sessionKeyPair"]["secretKey"],
            "wallet_public_key": data["session"]["walletPublicKey"],
            "bridge_url": data["session"]["bridgeUrl"]
        },
        "last_wallet_event_id": data["lastWalletEventId"],
        "connect_event": {
            "id": data["connectEvent"]["id"],
            "event": data["connectEvent"]["event"],
            "payload": {
                "items": [
                    {
                        "name": item["name"],
                        "address": item["address"],
                        "network": item["network"],
                        "publicKey": item["publicKey"],
                        "walletStateInit": item["walletStateInit"]
                    } for item in data["connectEvent"]["payload"]["items"]
                ],
                "device": data["connectEvent"]["payload"]["device"]
            }
        },
        "next_rpc_request_id": data["nextRpcRequestId"]
    }
    return json.dumps(session_data1)

# asyncio.run(save_session_data(chat_id,get_data(session_data)))
asyncio.run(main())
