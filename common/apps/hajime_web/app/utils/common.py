import hashlib
import uuid
from typing import TypeVar


import time
from datetime import datetime

from redis import asyncio as aioredis
from ulid import ULID

from app.db.schemas import GenericResponseModel
from datetime import datetime
import pytz

T = TypeVar('T')

def calculate_sha256(data):
    sha256_hash = hashlib.sha256()
    sha256_hash.update(data.encode('utf-8'))
    return sha256_hash.hexdigest()

def get_session_id()->str:
    random_uuid = uuid.uuid4()
    return calculate_sha256(str(random_uuid))
def get_unique_id():
    return str(ULID())
def get_current_time():
    return int(time.time() * 1000) + datetime.now().microsecond // 1000
def error_return(code, message, data:T=None):
    return GenericResponseModel(code=code, message= message,result=data)

def get_ts(date_string,format="%Y/%m/%d %H:%M:%S"):
    date_object = datetime.strptime(date_string, format)
    timestamp = date_object.timestamp()
    return timestamp


def get_time():
    # 创建一个时区对象
    tz = pytz.timezone('Asia/Shanghai')
    # 获取特定时区的当前时间
    current_time = datetime.now(tz)
    return current_time
def success_return(data: T=None):
    return GenericResponseModel(result=data)
async def get_redis():
    redis = await aioredis.from_url("redis://localhost",  db=0, decode_responses=True)
    return redis


def succ(data) -> T:
    data = {"code": 0, "message": "success", "result": data}
    return T(**data)