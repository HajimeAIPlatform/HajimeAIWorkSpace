# This is a sample Python script.
import asyncio
import json
import time
from datetime import datetime
import signal

import logging
from json import JSONDecodeError

import requests
from dotenv import load_dotenv

from app.database import init_db
from app.db.models import UserDeposit, LastCheckTime, MinerOrder, User, UserStat, OpLog
from app.db.schemas import UserDepositModel, SolScanModel
from app.utils.common import get_time

# 加载.env文件
load_dotenv()

api_key = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOjE3MTI3NzIwNzk2NjUsImVtYWlsIjoibm92YUBoYWppbWUuYWkiLCJhY3Rpb24iOiJ0b2tlbi1hcGkiLCJpYXQiOjE3MTI3NzIwNzl9.yIGBrunwuNwooWDnPen2nf9U-JKdnhdi-nY5JAwdT5E"


class SpiderTask:

    def __int__(self):
        pass

    async def get_transfer(self, account, offset=0, limit=20, fromTime=0):
        headers = {
            "token": api_key
        }
        url = f"https://pro-api.solscan.io/v1.0/account/splTransfers?account={account}&limit={limit}&offset={offset}"
        if fromTime > 0:
            now = int(time.time())
            # url = f"{url}&fromTime={fromTime}&toTime={now}"
        print(url)
        response = requests.get(url, headers=headers)
        new_num = 0
        try:
            data = response.json()
            # print(data)
            scan = SolScanModel(**data)
            # print("Total:", scan.total)
            for item in scan.data:
                exists = await  UserDeposit.find_one({"signature":{"$all":item.signature}})
                if exists is not None:
                    pass
                else:
                    new_num += 1
                    deposit = UserDeposit(**item.model_dump())
                    deposit.deleted = False
                    if deposit.address == "BfWYgztHDqrvnf1RXGofDR49JPi7BHbkCmdGDyqSHtKe":
                        deposit.nft_type = "D"  # 1499
                    elif deposit.address == "Gg8avohYTZ9G4skbqXJbtqzc99NmkTJfAhSomELWQSsh":
                        deposit.nft_type = "C"  # 20000
                    elif deposit.address == "9npi4xTUNBwWPKaCtcWDuVVR9zsrqmWRQqLzNkeQLujG":
                        deposit.nft_type = "A"  # 100000
                    elif deposit.address == "ALvmTrNzPuyJKzsHGZZBPhjCF39cXXJGhDzBwSguCGXF":
                        deposit.nft_type = "B"  # 50000
                    await deposit.save()
                    print(deposit)

            return new_num
        except:
            print("json_decode error", response.text)
            return 0

    async def get_sender(self, id, sign):
        """

        :param id: deposit id
        :param sign: tx_hash
        :return:
        """
        headers = {
            "token": api_key
        }
        url = f"https://pro-api.solscan.io/v1.0/transaction/{sign}"
        print(url)

        response = requests.get(url, headers=headers)
        try:
            data = response.json()
            if data['status'] == "Success":
                if len(data['signer']) == 1:
                    print(data['signer'], sign)
                    address = data['signer'][0]
                    sender = address.lower()
                    print(sender)
                    print(id)
                    # await UserDeposit.find_one({"_id": id}).update_one({"$set": {"sender": address}})


                else:
                    address = data['signer'][0]

                    sender = address.lower()

                    user = await User.get_user_by_address(sender)
                    if user is None:
                        await User(address=sender).create()
                    # print(address,token_address)
                    deposit = await UserDeposit.find_one({"_id": id})
                    amount = deposit.changeAmount / 1000000

                    user = await User.get_user_by_address(sender)
                    if user is not None:
                        exclude_list = ['nDeeLAYfCs6w7LiXQsoFQMyDnUmTk11YTHQ2mw1ytZy'.lower(),'fo3tPXfPU21FA2ixKgtL61DWtM6eeB86ozXERgRmwaR'.lower()]
                        if sender not in exclude_list:
                            await UserStat.add_score(user.id, amount)
                        else:
                            msg = "用户属于排出地址,不计数"
                            await OpLog.record(msg, user.id, sign, "exclude_user")


                        doc = {
                            "uid": user.id,
                            "pay_method": "crypto",
                            "status": 1,
                            "address": address.lower(),
                            "tx_hash": sign,
                            "amount": amount,
                            "nft_type": deposit.nft_type,
                        }
                        await MinerOrder(**doc).create()
                        await UserDeposit.find_one({"_id": id}).update_one({"$set": {"uid": user.id}})
                        msg = "扫描充值设置归属用户"
                        await OpLog.record(msg, user.id, sign, "set_deposit_uid")

                        await UserDeposit.find_one({"_id": id}).update_one({"$set": {"sender": address}})
        except JSONDecodeError:
            print("==========error=========")
            print(response)
            pass

    async def check_transfer(self):
        account_list = [
            "BfWYgztHDqrvnf1RXGofDR49JPi7BHbkCmdGDyqSHtKe",
            "Gg8avohYTZ9G4skbqXJbtqzc99NmkTJfAhSomELWQSsh",
            "9npi4xTUNBwWPKaCtcWDuVVR9zsrqmWRQqLzNkeQLujG",
            "ALvmTrNzPuyJKzsHGZZBPhjCF39cXXJGhDzBwSguCGXF"
        ]
        last_check_at = await LastCheckTime.get_last_check_time()

        for account in account_list:
            limit = 20
            offset = 0
            while True:
                num = await self.get_transfer(account, offset, limit, last_check_at)
                if num > 0:
                    offset += limit
                    print("offset,limit:", offset, limit)
                else:
                    print("done")
                    break

        await LastCheckTime.set_last_check_time()

    async def check_user_deposit(self, item: UserDeposit):
        if item.sender == "" and item.changeAmount > 0:
            signature = item.signature[0]
            await self.get_sender(item.id, signature)

    async def batch_check_user_deposit(self):
        items = await UserDeposit.find({"changeType": "inc", "changeAmount": {"$gt": 0}, "sender": ""}).to_list()
        for item in items:
            signature = item.signature[0]
            # print(signature)
            await self.get_sender(item.id, signature)
        items = await UserDeposit.find({"changeType": "dec", "sender": "", "changeAmount": {"$gt": 0}}).to_list()
        for item in items:
            signature = item.signature[0]
            # print(signature)
            await self.get_sender(item.id, signature)

    async def loop_work(self):

        while True:
            print("begin to check_transfer ",get_time())
            await self.check_transfer()
            print("begin to batch_check_user_deposit ",get_time())
            await self.batch_check_user_deposit()
            await asyncio.sleep(60)


async def main():
    # tid = await DBService.get_next("user_deposit_list")
    # print(tid)

    await init_db()
    executor = SpiderTask()

    # asyncio.create_task(scanner.fetch_price()),
    #
    # await scanner.scan_block(34439531)
    # return
    tasks = [
        asyncio.create_task(executor.loop_work()),

        # asyncio.create_task(executor.check_transfer()),
        # asyncio.create_task(executor.batch_check_user_deposit()),

    ]

    await asyncio.gather(*tasks)

    def signal_handler(sig, frame):
        loop = asyncio.get_event_loop()
        tasks = asyncio.all_tasks(loop=loop)
        for task in tasks:
            print("task:", task)
            task.cancel()
        loop.stop()

    signal.signal(signal.SIGINT, signal_handler)


# Press the green button in the gutter to run the script.
"""
Bot: 4XaKMf3Ge231LcMTDuu5DhdJBdAgBadkcvy62g6UDgQL
Bot 收款账号：BfWYgztHDqrvnf1RXGofDR49JPi7BHbkCmdGDyqSHtKe

C: 6MUweira2TG7eYHJYQzZQkWufBLYULzWrxBQ5eg5dMwa
C收款账号:Gg8avohYTZ9G4skbqXJbtqzc99NmkTJfAhSomELWQSsh

B: EcQeu6A5WTJntod5RTReAuY2epn1bd2hQmrwGXxe7Taq
B收款账号：ALvmTrNzPuyJKzsHGZZBPhjCF39cXXJGhDzBwSguCGXF

A: EMPC2TaeVL1hRJC8imXMEfa7kHqS1Fj55MAEoEqSDVnF
A收款账号：9npi4xTUNBwWPKaCtcWDuVVR9zsrqmWRQqLzNkeQLujG
"""
if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("async io")
        pass
