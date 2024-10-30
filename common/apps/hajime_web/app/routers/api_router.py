from typing import List, Optional

import pymongo
import pymysql
from fastapi import APIRouter, Depends, Response

from app.db.models import UserDeposit, User, UserStat, UserNewsSubscribe, UserAssetList, MinerOrder
from app.db.schemas import SignLoginModel, GenericResponseModel, UserAuthResponseModel, ReferModel, NewsSubscribeModel, \
    UserRewardListItemModel, UserRewardListResponse, UserInviteListItemModel, GeneralPageRequest, \
    UserInviteListResponse, UserStatResponse, MinerOrderRequest, UserInfoModel, MinerListResponse, MinerListItemModel, \
    APIDepositListItemModel, APIDepositListResponseModel
from app.service.user_service import UserService
from app.utils.common import error_return, success_return

router = APIRouter(prefix="/api", tags=["api"])

from pydantic import BaseModel, Field

from app.utils.x_auth import get_uid_by_token


# @router.get("/total_deposit")
# async def total_deposit():
#     return await UserDeposit.get_total()


# db = pymysql.connect(host='localhost', user='root', passwd='samos1688', port=3306,
#                      database="hajime")
# print('连接成功！')
# cursor = db.cursor()
#
# @router.get("/stat")
# async def stat():
#     items = await UserDeposit.find({}).to_list()
#     outs = {}
#     for item in items:
#         sender = item.sender
#         if sender != "":
#             if sender in outs:
#                 outs[sender] += item.changeAmount
#             else:
#                 outs[sender] = item.changeAmount
#     for k in outs:
#         outs[k] = outs[k]/1000000
#     for address in outs:
#         sql = "select * from wallet where address = '%s' and p_address !='' " % address
#         cursor.execute(sql)
#         results = cursor.fetchall()
#         print(results)
#         for row in results:
#             print(row)
#         # print(address,outs[address])
#     return outs





@router.post("/auth/login_with_sign", response_model=GenericResponseModel[UserAuthResponseModel], summary="授权登录")
async def login_with_sign(form: SignLoginModel, user_service: UserService = Depends()):
    ret = await  user_service.login_with_sign(form)
    return ret


@router.post("/auth/logout", response_model=GenericResponseModel, summary="退出")
async def logout(uid: str = Depends(get_uid_by_token), user_service: UserService = Depends()):
    ret = await  user_service.logout(uid)
    return ret


@router.post("/account/news_subscribe", response_model=GenericResponseModel, summary="订阅")
async def subscribe(form: NewsSubscribeModel, uid: str = Depends(get_uid_by_token),
                    user_service: UserService = Depends()):
    email = form.email.lower()
    user = await UserNewsSubscribe.find_one({"raw_email": email})
    if user is not None:
        return error_return(1, "email already subscribed")
    else:
        doc = {
            "raw_email": email,
            "email": form.email,
            "uid": uid,
        }
        await UserNewsSubscribe(**doc).create()
        return success_return(doc)



@router.post("/account/userinfo", response_model=GenericResponseModel, summary="用户信息")
async def userinfo( uid: str = Depends(get_uid_by_token),
                    user_service: UserService = Depends()):
    user = await User.get_user_by_uid(uid)
    doc = {
        "uid": user.id,
        "code": user.code,
        "invite_link": user.get_invite_link(),
        "pid": user.pid,
        "parent": user.p_address
    }

    return success_return(doc)



def callback_user_invite_list_item(item):
    return UserInviteListItemModel(**item.dict())




async def callback_api_invite_list_item(item):
    return UserInviteListItemModel(**item.dict())

@router.post("/account/user_invite_list", response_model=GenericResponseModel[UserInviteListResponse],
             response_description="会员邀请列表")
async def user_invite_list(form: GeneralPageRequest,
                    uid: str = Depends(get_uid_by_token)):

    query = {"pid": uid}

    options = {
        "page": form.page,
        "pagesize": form.pagesize,
        "sort": [("_id", pymongo.DESCENDING)],
    }
    return await UserStat.get_page(query, options, callback_api_invite_list_item, True)


async def callback_api_reward_list_item(item):
    return UserRewardListItemModel(**item.doc)

@router.post("/account/user_reward_list", response_model=GenericResponseModel[UserRewardListResponse],
             response_description="会员奖励列表")
async def user_reward_list(form: GeneralPageRequest,
                    uid: str = Depends(get_uid_by_token)):

    query = {"uid": uid}

    options = {
        "page": form.page,
        "pagesize": form.pagesize,
        "sort": [("_id", pymongo.DESCENDING)],
    }
    return await UserAssetList.get_page(query, options, callback_api_reward_list_item, True)

@router.post("/account/user_stat", response_model=GenericResponseModel[UserStatResponse],
             response_description="用户统计")
async def user_stat(
                    uid: str = Depends(get_uid_by_token)):
    user_stat = await UserStat.find_one({"uid": uid})
    if user_stat is None:
        return error_return(1, "user not found")
    else:
        stat = UserStatResponse(**user_stat.dict())
        return success_return(stat)

@router.post("/account/add_refer", response_model=GenericResponseModel, summary="绑定邀请码")
async def add_refer(form: ReferModel, uid: str = Depends(get_uid_by_token), user_service: UserService = Depends()):
    user = await  User.get_user_by_uid(uid)
    if user.pid != "":
        return error_return("已绑定")
    else:
        ret = await user_service.add_refer(uid,form.code)
        return ret





@router.post("/account/userinfo", response_model=GenericResponseModel, summary="用户信息")
async def userinfo( uid: str = Depends(get_uid_by_token)):
    user = await User.find_one({"_id":uid})
    if user is None:
        return error_return(code=1, message="user not found")
    else :

        userinfo = UserInfoModel(**user.dict())
        return success_return(userinfo)




@router.post("/market/make_order", response_model=GenericResponseModel, summary="下单")
async def make_order(form: MinerOrderRequest, uid: str = Depends(get_uid_by_token), user_service: UserService = Depends()):
    miner = await MinerOrder.find_one({"tx_hash":form.tx_hash})
    if miner is None:
        doc = form.dict()
        doc['uid']=uid
        await MinerOrder(**doc).create()
        return success_return()
    else:
        return error_return(code=1, message="already submitted")



async def callback_api_deposit_list_item(item:UserDeposit):
    # id: str
    # address: str
    # tx_hash: str
    # pay_method: str
    # product_id: int
    # pay_status: int
    # create_at: int
    # update_at: int=0
    doc = {
        "id": item.id,
        "address": item.sender,

        "amount": item.changeAmount/1000000,
        "tx_hash": item.signature[0],
        "sender": item.sender,
        "nft_type": item.nft_type,

        "create_at": item.blockTime,
    }
    return MinerListItemModel(**doc)






@router.post("/market/myorders", response_model=GenericResponseModel[MinerListResponse], summary="我的订单")
async def myorders( form:GeneralPageRequest,uid: str = Depends(get_uid_by_token)):
    query = {"uid": uid}

    print("uid:",uid)

    options = {
        "page": form.page,
        "pagesize": form.pagesize,
        "sort": [("blockTime", pymongo.DESCENDING)],
    }
    return await UserDeposit.get_page(query, options, callback_api_deposit_list_item, True)



    # options = {
    #     "page": form.page,
    #     "pagesize": form.pagesize,
    #     "sort": [("_id", pymongo.DESCENDING)],
    # }
    # return await UserStat.get_page(query, options, callback_user_miner_orderr_list_item, True)

    # items = await MinerOrder.find({"uid":uid}).to_list()
    # if items is not None:
    #     outs = []
    #     for item in items:
    #         miner_item = MinerListItemModel(**item.dict())
    #         outs.append(miner_item)
    #     return success_return(MinerListResponse(total=0, total_page=0, list=outs))
    # else:
    #     return error_return(code=1, message="already submitted")
#
#
#
# @router.get("/add_refer")
# async def add_refer():
#     items = await UserDeposit.find({}).to_list()
#     outs = {}
#     for item in items:
#         sender = item.sender
#         if sender != "":
#             if sender in outs:
#                 outs[sender] += item.changeAmount
#             else:
#                 outs[sender] = item.changeAmount
#     for k in outs:
#         outs[k] = outs[k]/1000000
#
#     reward = {}
#     output = []
#     for address in outs:
#
#         cursor = db.cursor()
#         sql = "select p_address,code from wallet where address = '%s' and p_address !='' " % address.lower()
#         # print(sql)
#         cursor.execute(sql)
#         results = cursor.fetchone()
#         # print(results)
#
#         if results is not None:
#             p_address = results[0]
#             if p_address in reward:
#                 reward[p_address] += outs[address]
#             else:
#                 reward[p_address] = outs[address]
#
#     for p_address in reward:
#         cursor = db.cursor()
#         sql = "select code from wallet where address = '%s'" % p_address.lower()
#         cursor.execute(sql)
#         results = cursor.fetchone()
#         code = results[0]
#         line = p_address + " " + str(reward[p_address]) + " " + code
#         output.append(line)
#
#         # print(address,outs[address])
#     return output

@router.post("/market/can_buy_miner", response_model=GenericResponseModel, summary="可以购买框架吗？")
async def can_buy_miner( ):
    data = {
        "can_buy_miner":0
    }
    return success_return(data)