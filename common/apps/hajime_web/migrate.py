# This is a sample Python script.
import asyncio
import json
import time
from datetime import datetime
import signal

import logging

import pymysql
import requests
from dotenv import load_dotenv
from pydantic import BaseModel

from app.database import init_db
from app.db.models import UserDeposit, LastCheckTime, User, UserStat
from app.db.schemas import UserDepositModel, SolScanModel
from app.service.user_service import UserService

# 加载.env文件
load_dotenv()

api_key = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOjE3MTI3NzIwNzk2NjUsImVtYWlsIjoibm92YUBoYWppbWUuYWkiLCJhY3Rpb24iOiJ0b2tlbi1hcGkiLCJpYXQiOjE3MTI3NzIwNzl9.yIGBrunwuNwooWDnPen2nf9U-JKdnhdi-nY5JAwdT5E"
db = pymysql.connect(host='localhost', user='root', passwd='samos1688', port=3306,
                     database="hajime")
print('连接成功！')
cursor = db.cursor()


class WalletModel(BaseModel):
    address: str
    p_address: str = ""
    code: str


class SpiderTask:

    def __int__(self):
        pass


    def get_wallet_items(self):
        cursor = db.cursor()
        cursor.execute("SELECT address,p_address,code FROM wallet")
        items = cursor.fetchall();
        outs = []
        for item in items:
            print(item)
            wallet = WalletModel(address=item[0], p_address=item[1], code=item[2])
            outs.append(wallet)
        return outs

    async def migrate_user(self):
        items = self.get_wallet_items()
        for item in items:
            # print(item)
            address = item.address.lower()


            user = await User.get_user_by_address(address)
            if user is None:
                user = await User(address=address, code=item.code).create()
                print("save user:", user)
            else:
                print("user exist:", address)

        items = self.get_wallet_items()
        user_server = UserService()
        for item in items:
            item.p_address = item.p_address.lower()
            address = item.address.lower()
            parent_user = await User.get_user_by_address(item.p_address)
            user = await User.get_user_by_address(address)
            if parent_user is not None and user is not None:
                await user_server.add_refer(user.id, parent_user.code)

    async def migrate_deposit(self):
        items = await UserDeposit.find({ "sender": {"$ne":""}}).to_list()
        for item in items:
            sender = item.sender.lower()
            if item.address == "BfWYgztHDqrvnf1RXGofDR49JPi7BHbkCmdGDyqSHtKe":
                item.nft_type = "A"  # 1499
            elif item.address == "Gg8avohYTZ9G4skbqXJbtqzc99NmkTJfAhSomELWQSsh":
                item.nft_type = "B"  # 20000
            elif item.address == "9npi4xTUNBwWPKaCtcWDuVVR9zsrqmWRQqLzNkeQLujG":
                item.nft_type = "D"  # 100000
            elif item.address == "ALvmTrNzPuyJKzsHGZZBPhjCF39cXXJGhDzBwSguCGXF":
                item.nft_type = "C"  # 50000
            await item.save()


            user = await User.get_user_by_address(sender)
            if user is not None:
                await UserStat.add_score(user.id, item.changeAmount)
            else:
                await User(address=sender).create()
                print("user not exist:,create user", sender)




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
        # asyncio.create_task(scanner.fetch_price()),
        asyncio.create_task(executor.migrate_user()),
        # asyncio.create_task(executor.migrate_deposit()),

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
