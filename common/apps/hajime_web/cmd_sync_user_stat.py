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





    async def check(self):
        items = await UserStat.find({"pid": ""}).to_list()
        total = 0
        for item in items:
            total = total + item.total_buy_amount
            print(item.uid,item.total_buy_amount)
        print(total)

    async def sync_counter(self):
        user_service = UserService()
        items = await UserStat.find().to_list()
        for item in items:
            user_service.clear_guard()
            data = await user_service.get_childs(item.uid)
            print(data)
            await UserStat.find_one({"uid":item.uid}).update_one({"$set":{"total_team_num":data['num']}})


        #
        # items = await UserStat.find().to_list()
        # for item in items:
        #     num = await  UserStat.find({"pid":item.uid}).count()
        #     await UserStat.find_one({"uid":item.uid}).update_one({"$set":{"direct_child":num}})
        #     print("child num:",num)


    async def sync_user_stat(self):
        # await UserStat.find({}).update_many({"$set":{""}})
        num1 = await User.count()
        num2 = await UserStat.count()
        print("num1:",num1)
        print("num2:",num2)
        users = await User.find({}).to_list()
        for user in users:
            uid = user.id
            stat = await UserStat.find_one({"uid":uid})
            if user.pid != stat.pid:
                if user.pid == "" and stat.pid != "":
                    user.pid = stat.pid
                    await  user.save()

                print("user_pid:",user.pid)
                print("stat_pid:",stat.pid)



        r = await UserStat.find({}).update_many({"$set": {"buy_amount":0,"total_buy_amount":0 }})
        print(r)
        items = await UserDeposit.find({"deleted":False}).to_list()
        for deposit in items:
            uid = deposit.uid
            amount = deposit.changeAmount / 1000000
            await UserStat.add_score(uid, amount)




        #
        #
        # uid = await User.get_uid_by_address("9748yjzsv6cxb8gnvatdkp7g3n9zt6vjfbtvj5xghpjb")
        # pid = await User.get_uid_by_address("bhaenqhobuvupdv6hvjngfoajktgktsnn2vrwseptgzf")
        #
        #
        #
        # user = await UserStat.find_one({"uid":uid})
        # parent = await UserStat.find_one({"uid":pid})
        # if user is not None and user.pid == "" and parent is not None :
        #     await UserStat.set_parent(uid,pid)



    async def clear(self):
        pass





# Set the signal handler

async def main():
    scanner = BscScanner()
    await init_db();

    tasks = [
        asyncio.create_task(scanner.sync_counter()),

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

