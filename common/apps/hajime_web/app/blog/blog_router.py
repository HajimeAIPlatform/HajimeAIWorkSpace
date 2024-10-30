import asyncio
from typing import List, Optional

import pymongo
import pymysql
from fastapi import APIRouter, Depends, Response

from app.blog.blog_model import Material, EmailSubscribe, Tag
from app.blog.blog_schema import MaterialModel, \
    AdminMaterialListRequestModel, AdminMaterialListItemModel, EmailSubscribeModel, AdminTagListItemModel, \
    PublishTypeModel
from app.db.models import User
from app.db.schemas import GenericResponseModel, IDModel
from app.utils.common import success_return, error_return, get_redis
from app.utils.x_auth import get_biz_uid_by_token

router = APIRouter(prefix="/blog", tags=["blog"])

from pydantic import BaseModel, Field


@router.post("/material_upsert", response_model=GenericResponseModel, response_description="添加修改素材")
async def material_upsert(form: MaterialModel, oid: int = Depends(get_biz_uid_by_token)):
    return  await Material.upsert(form)



async def callback_admin_material_list_item(item):
    doc = item.dict()
    return AdminMaterialListItemModel(**doc)

@router.post("/material_list", response_model=GenericResponseModel, response_description="列表")
async def material_list(form: AdminMaterialListRequestModel):
    query= {"is_published":1}

    if form.tag != "":
        query = {"tags":{"$in":[form.tag]}}
    # # print(query)
    # # print (request_form)

    options = {
        "page": form.page,
        "pagesize": form.pagesize,
        "sort": [("publish_at", pymongo.DESCENDING)],
    }
    #


    return await Material.get_page(query, options, callback_admin_material_list_item, True)


@router.post("/material_admin_list", response_model=GenericResponseModel, response_description="列表")
async def material_admin_list(form: AdminMaterialListRequestModel, oid: int = Depends(get_biz_uid_by_token)):
    query= {}

    if form.tag != "":
        query = {"tags":{"$in":[form.tag]}}
    # # print(query)
    # # print (request_form)

    options = {
        "page": form.page,
        "pagesize": form.pagesize,
        "sort": [("publish_at", pymongo.DESCENDING)],
    }
    #


    return await Material.get_page(query, options, callback_admin_material_list_item, True)


@router.post("/material_detail", response_model=GenericResponseModel, response_description="查看素材详情")

async def material_detail(form:IDModel):
    o = await Material.find_one({"_id":form.id})
    return success_return(MaterialModel(**o.model_dump()))


# async def material_del(form:IDModel  , oid: int = Depends(get_biz_uid_by_token)):

@router.post("/material_del", response_model=GenericResponseModel, response_description="删除素材详情")
async def material_del(form:IDModel  , oid: int = Depends(get_biz_uid_by_token)):
    r = await Material.find_one({"_id":form.id}).delete_one()
    if r.deleted_count == 1:
        return success_return("")
    else:
        return  error_return(404,"fail")



@router.post("/subscrie", response_model=GenericResponseModel, response_description="EMAIL订阅")
async def subscrie(form: EmailSubscribeModel):
    return  await EmailSubscribe.upsert(form)


@router.post("/tag_list", response_model=GenericResponseModel, response_description="Tag列表")
async def tag_list():
    query = {}
    items= await Tag.distinct("tag")
    print(items)
    # items = ['product','develop']
    return success_return(items)



async def callback_admin_tag_list_item(item):
    doc = item.dict()
    return AdminTagListItemModel(**doc)

@router.post("/tags", response_model=GenericResponseModel, response_description="列表")
async def tags(form: AdminMaterialListRequestModel):
    query = {}

    if form.tag != "":
        query = {"tags":{"$in":[form.tag]}}
    # # print(query)
    # # print (request_form)
    options = {
        "page": form.page,
        "pagesize": form.pagesize,
        "sort": [("publish_at", pymongo.DESCENDING)],
    }
    #


    return await Tag.get_page(query, options, callback_admin_tag_list_item, True)

@router.post("/tag_del", response_model=GenericResponseModel, response_description="删除素材详情")
async def tag_del(form:IDModel  , oid: int = Depends(get_biz_uid_by_token)):
    r = await Tag.find_one({"_id":form.id}).delete_one()
    if r.deleted_count == 1:
        return success_return("")
    else:
        return  error_return(404,"fail")


@router.post("/material_publish", response_model=GenericResponseModel, response_description="material_publish")
async def material_publish(form: IDModel, oid: int = Depends(get_biz_uid_by_token)):
    print(form)
    await Material.find_one({"_id":form.id}).update_one({"$set":{"is_published":1}})
    return  success_return()

@router.post("/material_unpublish", response_model=GenericResponseModel, response_description="material_unpublish")
async def material_publish(form: IDModel, oid: int = Depends(get_biz_uid_by_token)):
    print(form)

    await Material.find_one({"_id":form.id}).update_one({"$set":{"is_published":0}})
    return  success_return()

@router.post("/publish_www", response_model=GenericResponseModel, response_description="publish_www")
async def publish_www( oid: int = Depends(get_biz_uid_by_token)):
    redis = await get_redis()
    key = "publish"
    v = await redis.get(key)
    if v is not None:
        return error_return(1,"发布中")
    await redis.setex(key,90,1)
    import subprocess

    # 脚本的路径
    script_path = '/home/ubuntu/hajime-offcial-website/www.sh'
    print("script_path",script_path)

    # 调用脚本，这里假设deploy.sh需要一些参数
    try:
        # 使用 asyncio 创建一个异步的子进程
        proc = await asyncio.create_subprocess_exec(
            'bash', script_path, 'arg1', 'arg2',
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE
        )

        # 等待子进程完成
        stdout, stderr = await proc.communicate()

        # 打印输出
        print('STDOUT:', stdout.decode())
        print('STDERR:', stderr.decode())

        data = {
            "stdout": stdout.decode(),
            "stderr": stderr.decode()
        }
        await redis.delete(key)
        return success_return(data)

    except Exception as e:
        print('An error occurred while running the script:', e)
        await redis.delete(key)

        return error_return(1, str(e))

@router.post("/publish_wwwtest", response_model=GenericResponseModel, response_description="publish_wwwtest")
async def publish_wwwtest( oid: int = Depends(get_biz_uid_by_token)):
    redis = await get_redis()
    key = "publish"
    v = await redis.get(key)
    if v is not None:
        return error_return(1,"发布中")
    await redis.setex(key,90,1)

    # 脚本的路径
    script_path = '/home/ubuntu/hajime-offcial-website/wwwtest.sh'

    # 调用脚本，这里假设deploy.sh需要一些参数
    try:
        # 使用 asyncio 创建一个异步的子进程
        proc = await asyncio.create_subprocess_exec(
            'bash', script_path, 'arg1', 'arg2',
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE
        )

        # 等待子进程完成
        stdout, stderr = await proc.communicate()

        # 打印输出
        print('STDOUT:', stdout.decode())
        print('STDERR:', stderr.decode())

        data = {
            "stdout": stdout.decode(),
            "stderr": stderr.decode()
        }
        await redis.delete(key)


        return success_return(data)

    except Exception as e:
        print('An error occurred while running the script:', e)
        await redis.delete(key)

        return error_return(1, str(e))