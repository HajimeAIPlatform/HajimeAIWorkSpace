import logging
from typing import List




from fastapi import APIRouter, HTTPException, Depends

from app.db.schemas import GenericResponseModel, BizLoginModel
from app.service.user_service import UserService
from app.utils.common import get_redis
from app.utils.x_auth import get_biz_uid_by_token

from fastapi import Depends,Request
from fastapi.encoders import jsonable_encoder
from fastapi import FastAPI, Query
router = APIRouter(prefix="/admin", tags=["admin"])





@router.post("/biz_login", response_model=GenericResponseModel, response_description="登录后台")
async def biz_login(login_form: BizLoginModel,  request: Request,user_service: UserService = Depends()):
    print(login_form)
    return await user_service.biz_user_login(login_form)

@router.post("/biz_logout", response_model=GenericResponseModel, response_description="退出后台")
async def biz_logout(token: str = Depends(get_biz_uid_by_token),user_service: UserService = Depends()):
        if token:
            redis = await get_redis()
            oid = await redis.get(token)
            if oid is not None:
                await user_service.biz_logout(oid)
                await redis.delete(token)
                await redis.delete(oid)
            else:
                logging.error("logout error,no uid")
        else:
            logging.error("logout error,no token")
        return GenericResponseModel()