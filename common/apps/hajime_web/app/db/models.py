import asyncio
import base64
import json
import logging
import secrets
from decimal import Decimal
from hashlib import sha256
from typing import Any, List, Optional

import math

import pymongo
from beanie import Document, DecimalAnnotation, Insert, before_event, after_event, Indexed

import datetime

from pymongo import IndexModel

from app.config import CONFIG
from app.connection import get_next_id
from app.db.schemas import (GenericResponseModel, EvidenceDataInputModel, UserStateModel, AdminBuyGroup)
from app.utils.common import get_unique_id, get_current_time, error_return, success_return
from pydantic import  Field
import time

class BaseDocument(Document):

    @classmethod
    async def get_page(cls, query=None, options={"page": 1, "pagesize": 10, "sort": {"_id": -1}}, callback=None,
                       is_async_callback=False):
        if query is None:
            query = {}
        if options['page'] < 1:
            options['page'] = 1
        if options['pagesize'] < 10:
            options['pagesize'] = 10

        skip = (options['page'] - 1) * options['pagesize']
        limit = options['pagesize']
        sort = options['sort']

        total = await cls.find(query).count()
        total_page = math.ceil(total / options['pagesize'])

        objects = await cls.find(query).sort(sort).skip(skip).limit(limit).to_list()
        if callback is not None:
            if is_async_callback:
                objects = await asyncio.gather(*[callback(obj) for obj in objects])
            else:
                objects = list(map(callback, objects))

        out = {
            "total": total,
            "total_page": total_page,
            "pagesize": options['pagesize'],
            "list": objects
        }
        return GenericResponseModel(result=out)


class EvidenceData(BaseDocument):
    id: str = Field(default_factory=get_unique_id)
    task_id: str
    node_id:str=""
    data_hash: str
    raw_data: str = ""
    transaction_hash: str = ""
    data:Any
    status: int = 0
    create_at: int = get_current_time()
    update_at: int = get_current_time()

    class Settings:
        name = "evidence_data"

    @classmethod
    async def try_add_evidence(cls,form:EvidenceDataInputModel):
        try:
            data_hash = sha256(form.data.encode()).hexdigest()

            evidence = await cls.find_one({"task_id":form.task_id})
            if evidence is not None:
                return error_return(code=1, message="task_id exists")
            decoded_bytes = base64.b64decode(form.data)
            decoded_str = decoded_bytes.decode('utf-8')
            data = json.loads(decoded_str)
            doc = {
                "data_hash": data_hash,  #
                "raw_data": form.data,
                "task_id": form.task_id,
                "data": data,
                "node_id":form.node_id
            }
            evidence = await EvidenceData(**doc).create()

            return success_return({"task_id": evidence.task_id,"data_hash": evidence.data_hash,"node_id":form.node_id})

        except Exception as e:
            print(e)
            return error_return(code=500, message="data error",data=form)

    @classmethod
    async def update_task_hash(cls,task_id,hash):
        return await cls.find_one({"task_id":task_id}).update_one({"$set":{"transaction_hash":hash,"status":1,"update_at":get_current_time()}})


class UserDeposit(BaseDocument):
    id: str = Field(default_factory=get_unique_id)
    uid:str = ""
    address: str=""
    signature: List[str]
    changeType: str
    changeAmount: int
    decimals: int
    postBalance: int
    preBalance: int
    tokenAddress: str
    owner: str
    blockTime: int
    slot: int
    fee: int
    symbol: str
    sender:str=""
    nft_type:str=""
    tx_hash:str=""
    deleted:bool=False
    create_at: int = get_current_time()
    update_at: int = get_current_time()

    class Settings:
        name = "user_deposit"
        indexes = [
            IndexModel(
                [("signature", pymongo.DESCENDING)],
                name="signature_list_index_DESCENDING",
            )
        ]

    @classmethod
    async def get_user_total(cls,uid):
        total = await cls.find({"changeAmount":{"$gte":0},"uid":uid,"deleted":False}).sum("changeAmount")
        # print(total)
        if total is not None:
            return total/1000000
        else:
            return 0
    @classmethod
    async def get_total(cls):
        total1 = await cls.find({"changeAmount": {"$gte": 0}}).sum("changeAmount")
        print(total1)

        total = await cls.find({"changeAmount":{"$gte":0},"deleted":False}).sum("changeAmount")
        print(total)
        if total is not None:
            return total/1000000
        else:
            return 0

    @classmethod
    async def get_buy_group(cls):
        items = await cls.find(
            {"changeAmount": {"$gte": 0},"deleted":False}).aggregate(
            [{"$group": {"_id": "$nft_type", "total": {"$sum": "$changeAmount"}}}],
            projection_model=AdminBuyGroup
        ).to_list()
        outs = []
        for item in items:
            o = item
            o.total = item.total/1000000
            outs.append(o)
        return outs
    @classmethod
    async def get_total_by_query(cls,query):

        total = await cls.find(query).sum("changeAmount")
        # print(total)
        if total is not None:
            return total/1000000
        else:
            return 0

class LastCheckTime(BaseDocument):
    id: str = Field(default_factory=get_unique_id)
    last_check_at: int
    class Settings:
        name = "last_check_time"


    @classmethod
    async def get_last_check_time(cls):
        vo = await cls.find_one()
        if vo:
            return vo.last_check_at
        else:
            return 0

    @classmethod
    async def set_last_check_time(cls):
        vo = await cls.find_one()
        if vo:
            vo.last_check_at = int(time.time())-120
            await vo.save()
        else:
            doc = {
                "last_check_at": int(time.time())-120
            }
            await LastCheckTime(**doc).create()


class MinerOrder(BaseDocument):
    id: str = Field(default_factory=get_unique_id)
    uid:str=""
    pay_method: str="crypto"
    address:str=""
    product_id:int=1
    tx_hash:str=""
    status:int=0
    price:str=""
    create_at: int = get_current_time()
    update_at: int = get_current_time()
    class Settings:
        name = "miner_order"


class UserStat(BaseDocument):
    id:str=Field(default_factory=get_unique_id)
    uid: str
    address: Optional[str] = ""  #
    pid: str = ""
    p_address: Optional[str] = ""  #
    split_pid:str = ""
    buy_amount:  DecimalAnnotation = Field(decimal_places=8, default=0 )
    total_buy_amount:  DecimalAnnotation = Field(decimal_places=8, default=0 )
    direct_active_child:int=0
    direct_child:int=0
    total_active_team_num:int=0
    total_team_num:int=0
    active: int = 0
    create_at: int = get_current_time()
    update_at: int = get_current_time()


    @classmethod
    async def get_active_child_list(cls,pid):
        items  = await cls.find({"pid":pid}).to_list()
        return items

    @classmethod
    async def get_active_parent(cls,pid):
        user = await cls.find_one({"uid":pid})
        if user and user.active:
            return user.pid
        else:
            pid = user.pid
            i = 1
            while i < 10000:
                print("i=>",i)
                user = await cls.find_one({"uid": pid})
                if user is None:
                    return None
                if user.active:
                    return user.pid
                else:
                    pid = user.pid
                i+=1

    @classmethod
    async def set_split_parent(cls,pid,uid):
        await cls.find_one({"uid":uid}).update({"$set":{"split_pid":pid}})

    @classmethod
    async def set_parent(cls,uid,pid):
        user  = await cls.find_one({"uid":uid})
        parent = await cls.find_one({"uid":pid})
        if user is not None and parent is not None and user.pid == "":
            r1 = await cls.find_one({"uid":uid}).update({"$set":{"pid":pid,"p_address":parent.address}})
            r2 = await cls.find_one({"uid":pid}).update({"$inc":{"direct_child":1}})
            await User.find_one({"uid":uid}).update({"$set":{"pid":pid,"p_address":parent.address}})
            print(r1)
            print(r2)

            pid_list = await cls.get_parents(uid)
            r3 = await cls.find({"uid":{"$in":pid_list}}).update_many({"$inc":{"total_buy_amount":user.total_buy_amount,"total_team_num":1}})
            print(r3)



    @classmethod
    async def get_parents(cls,uid,height=1000):
        ids = []
        for i in range(0,height):
            user = await cls.find_one({"uid":uid}).project(UserStateModel)
            if user is None:
                break
            else:
                if user.pid not in ids:
                    ids.append(user.pid)
                uid = user.pid

        return ids

    @classmethod
    async def add_score(cls,uid, amount):
        await cls.find_one({"uid":uid}).update({"$inc":{"total_buy_amount":amount,"buy_amount":amount}})
        msg="增加购买"
        extra = {"amount": amount}
        await OpLog.record(msg,uid,extra,"add_score")
        id_set = await cls.get_parents(uid)
        if id_set:
            await cls.find({"uid":{"$in":id_set}}).update_many({"$inc":{"total_buy_amount":amount}})
            msg = "给父亲节点增加购买"
            await OpLog.record(msg, uid, id_set, "add_score")

    @classmethod
    async def dec_score(cls,uid, amount):
        if amount > 0:
            amount = amount * (-1)
        await cls.find_one({"uid":uid}).update({"$inc":{"total_buy_amount":amount,"buy_amount":amount}})
        msg="减少购买"
        extra = {"amount": amount}
        await OpLog.record(msg,uid,extra,"dec_score")
        id_set = await cls.get_parents(uid)
        if id_set:
            await cls.find({"uid":{"$in":id_set}}).update_many({"$inc":{"total_buy_amount":amount}})
            msg = "给父亲节点减少购买"
            await OpLog.record(msg, uid, id_set, "dec_score")

    class Settings:
        name = "user_stat"
        indexes = [
            [
                ("uid", pymongo.ASCENDING),
                ("pid", pymongo.ASCENDING)
            ],
        ]
    class Config:
        json_schema_extra = {
            "example": {
                "uid": "1",
                "username": "bots001",
                "address": "0xE91C299427D5DB24Dcc064db3Dc42d1bF1bf187E",

            }
        }

class UserAssetList(BaseDocument):
    id:str=Field(default_factory=get_unique_id)
    uid: str=""
    mainchain: str=""
    token: str=""
    amount:  DecimalAnnotation = Field(decimal_places=8, default=0 )
    op_type: str = ""
    desc: str = ""
    status:str ="unconfirmed"
    create_at: int = get_current_time()
    update_at: int = get_current_time()
    class Settings:
        name = "user_asset_list"
        indexes = [
            [
                ("uid", pymongo.ASCENDING),
                ("op_type", pymongo.ASCENDING)
            ],
        ]
    class Config:
        json_schema_extra = {
            "example": {
                "uid": "1",
                "username": "bots001",
                "address": "0xE91C299427D5DB24Dcc064db3Dc42d1bF1bf187E",
            }
        }

class UserWithdraw(Document):
    id: str = Field(default_factory=get_unique_id)
    uid: str
    withdraw_id:int=0
    mainchain: str
    token: str
    address: str
    amount:  DecimalAnnotation = Field(decimal_places=8, default=0.00000000 )
    status: Indexed(str)
    check_times:int=0
    hash: Indexed(str)=""
    src_id:str=""
    create_at: int = get_current_time()
    update_at: int = get_current_time()


    class Settings:
        name = "user_withdraw"
        indexes = [
            [
                ("address", pymongo.TEXT),
                ("uid", pymongo.ASCENDING)
            ],
        ]

    class Config:
        populate_by_name = True
        arbitrary_types_allowed = True
        json_schema_extra = {
            "example": {
                "uid": "1111",
                "address": "0xE91C299427D5DB24Dcc064db3Dc42d1bF1bf187E",
            }
        }

    @classmethod
    async def add_withdraw(cls,uid,address,mainchain,token,amount,desc="",status="waiting"):
        doc = {
            "withdraw_id":await get_next_id("withdraw"),
            "uid":uid,
            "mainchain":mainchain,
            "token":token,
            "amount":amount,
            "address":address,
            "desc":desc,
            "status":status
        }
        await cls(**doc).create()
        logging.info("提现:"+token+" ,"+amount+",{desc}",uid)


class UserAsset(BaseDocument):
    id: str = Field(default_factory=get_unique_id)
    uid: str
    mainchain: str = "SOLANA"
    token: str = "SOL"
    amount:  DecimalAnnotation = Field(decimal_places=8, default=0.00000000 )
    frozen:  DecimalAnnotation = Field(decimal_places=8, default=0.00000000 )

    create_at: int = get_current_time()
    update_at: int = get_current_time()

    @classmethod
    async def init_token(cls, uid, mainchain, token):
        vo = await cls.find_one({"uid": uid, 'mainchain': mainchain, "token": token})
        if vo:
            return vo
        else:
            item = await  UserAsset(id=get_unique_id(), uid=uid, mainchain=mainchain, token=token).save()
            return item

    def get_amount(self):
        return str((self.amount))

    def get_frozen(self):
        return str((self.frozen))

    @classmethod
    async def get_user_asset_list(cls, uid):
        return await cls.find({"uid": uid}).to_list()

    @classmethod
    async def get_user_asset_map(cls, uid):
        out = {}
        items = await cls.find({"uid": uid}).to_list()
        for item in items:
            out[item.token] = (item.amount)
        return out

    @classmethod
    async def get_user_filter_asset_map(cls, where):
        out = {}
        items = await cls.find(where).to_list()
        for item in items:
            out[item.token] = ((item.amount))
        return out

    @classmethod
    async def incr(cls, uid: str, mainchain: str, token: str, amount: Decimal, op_type: str = "", desc: str = ""):
        vo = await cls.find_one({"uid": uid, "token": token,"mainchain":mainchain})
        if vo is None:
            await cls.init_token(uid, mainchain, token)
        if amount < 0:
            r = await cls.find_one({"uid": uid, "mainchain": mainchain,
                                    "token": token, "amount": {"$gte": amount * (-1)}}).update(
                {"$inc": {"amount": amount}})
        else:
            r = await cls.find_one({"uid": uid, "mainchain": mainchain, "token": token}).update(
                {"$inc": {"amount": amount}})
        i = r.modified_count
        print("incr->reuslt:", i,token,amount)

        if r.modified_count > 0:
            doc = {
                "uid": uid,
                "mainchain": mainchain,
                "token": token,
                "amount": amount,
                "op_type": op_type,
                "desc": desc
            }
            vo = await UserAssetList(**doc).create()
            return success_return({"id":vo.id})
        else:
            return error_return(1,"insufficent_balance",{"amount":amount})

    @classmethod
    async def withdraw(cls,uid,address,token,num,op_type="withdraw",desc="withdraw",mainchain="BSC"):
        print(token,num)
        ret = await cls.incr( uid, mainchain, token, Decimal(num) * (-1) , op_type, desc)
        if ret.code == 0:
            await UserWithdraw.add_withdraw(uid,address,mainchain,token,num,desc="user withdraw")
            return success_return(num)
        else:
            return error_return(1,ret.message)

    class Settings:
        name = "user_asset"
        indexes = [
            [
                ("address", pymongo.TEXT),
                ("uid", pymongo.ASCENDING)
            ],
        ]

    class Config:
        populate_by_name = True
        arbitrary_types_allowed = True
        json_schema_extra = {
            "example": {
                "uid": "1111",
                "address": "0xE91C299427D5DB24Dcc064db3Dc42d1bF1bf187E",
            }
        }



class User(BaseDocument):
    id: str = Field(default_factory=get_unique_id)
    pid: Optional[str] = ""
    address: Optional[str] = ""  #
    p_address: Optional[str] = ""  #
    email: Optional[str] = ""
    avatar: Optional[str] = ""
    password: Optional[str] = ""
    trade_password: Optional[str] = ""
    twitter: Optional[str] = ""
    telegram: Optional[str] = ""
    discord: Optional[str] = ""
    status:int=1
    code :str = ""
    create_at: int = get_current_time()
    update_at: int = get_current_time()

    class Settings:
        name = "user"


    @classmethod
    async def get_uid_by_address(cls,addrress):
        vo =  await cls.find_one({"address": addrress.lower()})
        if vo is not None:
            return vo.id
        else:
            return ""
    @classmethod
    async def get_address_by_uid(cls,uid):
        vo =  await cls.find_one({"_id": uid})
        if vo is not None:
            return vo.address
        else:
            return ""
    @classmethod
    async def get_user_by_uid(cls,uid):
        return await cls.find_one({"_id": uid})

    @classmethod
    async def get_user_by_address(cls, address:str):
        return await cls.find_one({"address": address.lower()})
    @classmethod
    async def get_user_by_code(cls, code:str):
        return await cls.find_one({"code": code})
    @classmethod
    async def set_parent(cls, uid, pid):
        await cls.find_one(User.id == uid).update({"$set": {"pid": pid}})
        await UserStat.find_one(User.id == uid).update({"$set": {"pid": pid}})


    @classmethod
    async def get_parents_ids(cls, uid, height=10000):
        parents = []
        i = 0
        while i < height:
            user = await cls.find_one(User.id == uid)
            if user:
                if i != 0:  # not self
                    parents.append(user.id)
                if user.pid != "":  # have parent
                    uid = user.pid
                else:
                    break
            else:
                break
            i = i + 1

        return parents


    def get_invite_link(self):
        return CONFIG.invite_domain + "#/home?code=" + self.code

    @before_event(Insert)
    async def generate_invite_code(self):
        if not self.code:
            while True:
                code = secrets.token_hex(4)
                user = await User.find_one(User.code == code)
                if not user:
                    self.code = code
                    break
        else:
            print("exists code")


    @after_event(Insert)
    async def init_stat(self):

        self.create_at = get_current_time()
        await UserStat(uid=self.id, pid=self.pid, address=self.address).create()
        logging.info(("user stat created"))

        for token in CONFIG.token_list:
            await UserAsset(uid=self.id, mainchain="BSC", token=token).create()





class BizUser(BaseDocument):
    id:str=Field(default_factory=get_unique_id)
    password: str = Field(..., min_length=5, max_length=256)
    username: str = Field(..., min_length=5, max_length=40)
    ga: str=""
    avatar: str=""
    status:int=1
    create_at: int = Field(default_factory=get_current_time)
    update_at: int = Field(default_factory=get_current_time)

    class Settings:
        name = "biz_user"
        indexes = [
            IndexModel([("username", pymongo.ASCENDING)], unique=True)
        ]



    def get_avatar(self):
       if  self.avatar != "":
           return self.avatar
       else:
           return "http://xx.xx"


    class Config:
        json_schema_extra = {
            "example": {
                "username": 1,
                "password": "edu",
            }
        }



class UserNewsSubscribe(BaseDocument):
    id:str=Field(default_factory=get_unique_id)
    uid: str
    email:str
    raw_email:str
    create_at: int = Field(default_factory=get_current_time)
    update_at: int = Field(default_factory=get_current_time)

    class Settings:
        name = "user_news_subscribe"


class OpLog(BaseDocument):
    msg: str
    uid:str
    extra: str=""
    op:str
    create_at: int = Field(default_factory=get_current_time)

    @classmethod
    async def record(cls,msg,uid:int,extra={},op=""):
        data = {
            "uid":uid,
            "msg":msg,
            "op":op,
            "extra":json.dumps(extra),
        }
        await OpLog(**data).create()

    class Settings:
        name = "oplog"


