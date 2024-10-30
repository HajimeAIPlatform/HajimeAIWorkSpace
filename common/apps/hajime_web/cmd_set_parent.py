# This is a sample Python script.
import asyncio
import signal
import sys

import os


import logging
from decimal import Decimal

from dotenv import load_dotenv

from app.db.models import User, UserStat, UserAsset, MinerOrder, UserDeposit, LastCheckTime
from app.db.schemas import BizUserModel
from app.service.user_service import UserService

# 加载.env文件
load_dotenv()

# 读取变量


logging.basicConfig(filename='scan.log', level=logging.INFO)



from app.database import init_db





class BscScanner(object):


    def __init__(self):
        self.bnb_price = None






    async def set_parent(self):

        uid = await User.get_uid_by_address("9748yjzsv6cxb8gnvatdkp7g3n9zt6vjfbtvj5xghpjb")
        pid = await User.get_uid_by_address("bhaenqhobuvupdv6hvjngfoajktgktsnn2vrwseptgzf")



        user = await UserStat.find_one({"uid":uid})
        parent = await UserStat.find_one({"uid":pid})
        if user is not None and user.pid == "" and parent is not None :
            await UserStat.set_parent(uid,pid)



    async def clear(self):
        pass





# Set the signal handler

async def main():
    scanner = BscScanner()
    await init_db();

    tasks = [
        asyncio.create_task(scanner.set_parent()),

    ]
    await asyncio.gather(*tasks)



    def signal_handler(sig, frame):
        loop = asyncio.get_event_loop()
        tasks = asyncio.all_tasks(loop=loop)
        for task in tasks:
            print("task:",task)
            task.cancel()
        loop.stop()


    signal.signal(signal.SIGINT, signal_handler)

# Press the green button in the gutter to run the script.

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("async io")
        pass

