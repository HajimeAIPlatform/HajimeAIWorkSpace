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


    async def check_amount(self):
        items  = await UserStat.find({}).to_list()
        for item in items:
            uid = item.uid
            total = await UserDeposit.get_user_total(uid)
            if total > item.buy_amount:
                print("==>Total miss:",total,item.buy_amount)
            elif total < item.buy_amount:
                print("<==Total extra:",total,item.buy_amount)
                extra = float(item.buy_amount) - total
                amount = Decimal(extra)* (-1)
                ids = await UserStat.get_parents(uid)
                await UserStat.find_one({"uid":uid}).update_one({"$inc":{"buy_amount":amount ,"total_buy_amount":amount}})
                await UserStat.find({"uid":{"$in":ids}}).update_many({"$inc":{"total_buy_amount":amount}})
                print("sub uid:",uid,ids,amount)


    async def set_nft_type(self):
        items = await UserDeposit.find({}).to_list()
        for item in items:
            if item.changeAmount >0:
                if item.changeAmount == 1499000000:
                    await UserDeposit.find_one({"_id":item.id}).update_one({"$set":{"nft_type":"D"}})
                elif  item.changeAmount == 20000000000:
                    await UserDeposit.find_one({"_id":item.id}).update_one({"$set":{"nft_type":"C"}})
                elif  item.changeAmount == 50000000000:
                    await UserDeposit.find_one({"_id":item.id}).update_one({"$set":{"nft_type":"B"}})
                elif  item.changeAmount == 100000000000:
                    await UserDeposit.find_one({"_id":item.id}).update_one({"$set":{"nft_type":"A"}})




    async def set_parent(self):
        uid = "01HVP6AC3JTJVBW6HF952JXAGC"
        pid = "01HVP6ABZJSTQ5MDR4QCPHGCWA"
        user = await UserStat.find_one({"uid":uid})
        parent = await UserStat.find_one({"uid":pid})
        if user.pid == "":
            await UserStat.set_parent(uid,pid)


    async def dec_score(self):
        ids = [
"665b3c68f36d498a90794b79",

        ]
        for id in ids:
            print("query :",id)
            depoist = await UserDeposit.find_one({"_id":id})
            if depoist is not None and depoist.deleted == False:
                amount = depoist.changeAmount
                uid = depoist.uid
                if amount < 0:
                    amount = amount * (-1)
                amount = amount / 1000000
                print(depoist)
                await UserStat.dec_score(uid,amount)
                depoist.deleted = True
                await depoist.save()
                print("sleeping\n")
                await asyncio.sleep(6)
                print("sleeping done\n\n")


    async def add_oid(self):
        data = {
            "username":"blogeditor",
            "password" :"blogeditor168816999",
            "ga":""
        }
        form = BizUserModel(**data)
        user_service = UserService()
        await user_service.create_biz_user(form)
    async def do_task(self):
        await self.add_oid()
    async def clear(self):
        pass

    async def stat(self):
        total = 0
        outs = {}
        parent_map = {}
        query = {"changeAmount": {"$gt": 0}, "deleted": False}
        fp = open("c.csv","w")
        items = await UserDeposit.find(query).to_list()
        for item in items:
            sender = item.sender
            exclude_list = ['nDeeLAYfCs6w7LiXQsoFQMyDnUmTk11YTHQ2mw1ytZy'.lower(),
                            '3nDeeLAYfCs6w7LiXQsoFQMyDnUmTk11YTHQ2mw1ytZy'.lower(),
                            '9fo3tPXfPU21FA2ixKgtL61DWtM6eeB86ozXERgRmwaR'.lower(),
                            'fo3tPXfPU21FA2ixKgtL61DWtM6eeB86ozXERgRmwaR'.lower()]
            if sender.lower() in exclude_list:
                continue
            if sender not in outs.keys():
                outs[sender] = {}
            nft_type = item.nft_type
            uid = item.uid
            amount = int(item.changeAmount /1000000)
            if nft_type in ['A','B','C','D'] and amount>0:
                stat = await UserStat.find_one({"uid":uid})
                pid = stat.pid
                if len(pid) > 10:
                    parent_user = await User.get_user_by_uid(pid)
                    # print(stat)
                    parent = parent_user.address
                    total += amount
                    parent_map[sender] = parent


                    # parent = await User.get_address_by_uid(pid)
                    # print("uid:", uid)
                    # print("pid:", pid)
                    #
                    # print(sender, ",",nft_type,",", amount, ",", parent)
                    # print(("\n=====\n"))

                else:
                    parent = ""
                    # print(sender, ",", nft_type, ",", amount, ",", )
                    total += amount
                    pid = ""
                    parent = ""
                # print("uid:",uid)
                # print("pid:",pid)
                # print(nft_type,",",sender,",",amount,",",parent)
                # print(("\n=====\n"))
                if nft_type not in outs[sender].keys():
                    outs[sender][nft_type] = amount
                else:
                    outs[sender][nft_type] += amount

        # print("Total:",total)
        # print(outs)
        for sender in outs:
            parent = "None"
            if sender in parent_map.keys():
                parent = parent_map[sender]
            for nft_type in outs[sender]:
                if nft_type == 'D':
                    num = outs[sender][nft_type]/1499
                elif nft_type == 'C':
                    num = outs[sender][nft_type]/20000
                elif nft_type == 'B':
                    num = outs[sender][nft_type] / 50000
                elif nft_type == 'A':
                    num = outs[sender][nft_type] / 100000

                num = int(num)


                line = sender+","+nft_type+","+str(num)+","+str(outs[sender][nft_type])+","+parent+"\n"
                print(line)
                fp.write(line)

        # print("Total:",total)
        fp.close()


# Set the signal handler

async def main():
    scanner = BscScanner()
    await init_db();

    tasks = [
        asyncio.create_task(scanner.add_oid()),

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

