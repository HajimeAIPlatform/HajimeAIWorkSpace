from typing import List

import pymongo
from fastapi import APIRouter, HTTPException, Depends

from fastapi import Depends, Request
from fastapi.encoders import jsonable_encoder
from fastapi import FastAPI, Query
from pydantic import BaseModel

from app.db.models import User, UserStat, UserDeposit, OpLog
from app.db.schemas import GenericResponseModel, BizLoginModel, BizAuthResponseModel, UserListRequestModel, \
    UserListItemModel, AdminUserListItemModel, AdminUserListResponseModel, UIDModel, AdminDepositListItemModel, \
    AdminDepositListRequestModel, AdminDepositListResponseModel, AdminOplogListItemModel, AdminOplogListResponseModel, \
    AdminOplogListRequestModel, IDModel, UploadResponseModel
from app.service.fs_service import FsService
from app.service.user_service import UserService
from app.utils.common import success_return, get_ts, error_return
from app.utils.x_auth import get_uid_by_token, get_biz_uid_by_token
from fastapi import FastAPI, File, UploadFile,Request

router = APIRouter(prefix="/admin", tags=["admin"])


async def callback_admin_user_list_item(item: UserStat):
    doc = item.dict()
    if len(item.pid) > 10:
        parent_user = await User.get_user_by_uid(item.pid)
        doc["p_address"] = parent_user.address
    user = await User.get_user_by_uid(item.uid)
    if user is not None:
        doc["address"] = user.address
        doc["code"] = user.code
    return AdminUserListItemModel(**doc)


@router.post("/account/users", response_model=GenericResponseModel[AdminUserListResponseModel],
             response_description="会员列表")
async def admin_user_list(request_form: UserListRequestModel,
                          oid: int = Depends(get_biz_uid_by_token)):
    query = {}
    if len(request_form.address) > 0:
        user = await User.get_user_by_address(request_form.address.lower())
        if user is not None:
            query = {"uid": user.id}
    options = {
        "page": request_form.page,
        "pagesize": request_form.pagesize,
        "sort": [("total_team_num", pymongo.DESCENDING)],
    }
    return await UserStat.get_page(query, options, callback_admin_user_list_item, True)


@router.post("/stat/dashboard", response_model=GenericResponseModel,
             response_description="数据面板")
async def admin_dashboard(
                          oid: int = Depends(get_biz_uid_by_token)):

    members = await User.count()
    total_invested = await UserDeposit.get_total()
    doc = {
        "total_user": members,
        "total_invested":total_invested,
        "buy_group_list":await UserDeposit.get_buy_group()
    }

    return success_return(doc)


@router.post("/stat/getchilds", response_model=GenericResponseModel,
             response_description="获取子节点数据")
async def getchilds(form: UIDModel,
                    oid: int = Depends(get_biz_uid_by_token),
                    user_service: UserService = Depends()):
    data = await  user_service.get_all_childs(form.uid)

    user_service.clear_guard()
    out = {
        'nodes': data[0],
        'num': data[1]
    }
    return GenericResponseModel(code=0, message="success", result=out)


async def callback_deposit_list_item(item: UserDeposit):
    item.changeAmount = item.changeAmount / 10 ** item.decimals

    if len(item.sender)>3:
        item.uid = await User.get_uid_by_address(item.sender)
    doc = item.dict()
    doc['create_at'] = item.blockTime * 1000



    return AdminDepositListItemModel(**doc)


@router.post("/asset/depositelist", response_model=GenericResponseModel[AdminDepositListResponseModel],
             response_description="充值列表")
async def depositlist(request_form: AdminDepositListRequestModel,
                      oid: int = Depends(get_biz_uid_by_token)):
    query = {"changeAmount":{"$gt":0},"deleted":False}
    print(request_form)
    if len(request_form.username) > 0:
        user = await User.get_user_by_address(request_form.username.lower())
        if user is not None:
            query["uid"] = user.id
    if len(request_form.type) >0:
        type = request_form.type.upper()
        query['nft_type'] = type
    # if len(request_form.start_date)>3:
    #     query['create_at']={'$gte':get_ts(request_form.start_date)*1000}
    # if len(request_form.end_date) > 3:
    #     query['create_at']={'$lt':get_ts(request_form.end_date)*1000}
    if  len(request_form.start_date)>3 and len(request_form.end_date) > 3:

        query['$and']=[{'blockTime':{'$lt':get_ts(request_form.end_date)}},
                       {'blockTime':{'$gte':get_ts(request_form.start_date)}}]


    print(query)

    # print (request_form)
    options = {
        "page": request_form.page,
        "pagesize": request_form.pagesize,
        "sort": [("blockTime", pymongo.DESCENDING)],
    }

    total = await UserDeposit.get_total_by_query(query)
    ret =  await UserDeposit.get_page(query, options, callback_deposit_list_item, True)
    ret.result['total_deposit'] = total
    return ret



async def callback_oplog_list_item(item: OpLog):
    doc = item.dict()
    doc['address'] = await User.get_address_by_uid(item.uid)
    return AdminOplogListItemModel(**doc)


@router.post("/operation/oplog", response_model=GenericResponseModel[AdminOplogListResponseModel],
             response_description="运行日志")
async def oplog(request_form: AdminOplogListRequestModel,
                      oid: int = Depends(get_biz_uid_by_token)):
    query = {}
    if len(request_form.username) > 0:
        user = await User.get_user_by_address(request_form.username.lower())
        if user is not None:
            query["uid"] = user.id

    # print(query)
    # print (request_form)
    options = {
        "page": request_form.page,
        "pagesize": request_form.pagesize,
        "sort": [("blockTime", pymongo.DESCENDING)],
    }

    ret =  await OpLog.get_page(query, options, callback_oplog_list_item, True)
    return ret


@router.post("/admin_upload", response_model=GenericResponseModel[UploadResponseModel],tags=["admin"],summary="上传文件")
async def admin_upload(uid:str = Depends(get_biz_uid_by_token),
                             file: UploadFile = File(...),
                             file_service: FsService = Depends()):
    """
        上传图片/文件 私有的,私有上传的要用下面的下载接口下载
    """
    is_private = True
    data =  await file_service.deal_file(uid,file,is_private)
    return GenericResponseModel(result=data)
