import hashlib
import uuid
from typing import TypeVar
import time
from datetime import datetime
from ulid import ULID
import pytz
from pythonp.apps.hajime_blog.db.schemas import GenericResponseModel

T = TypeVar('T')

def calculate_sha256(data: str) -> str:
    sha256_hash = hashlib.sha256()
    sha256_hash.update(data.encode('utf-8'))
    return sha256_hash.hexdigest()

def get_session_id() -> str:
    random_uuid = uuid.uuid4()
    return calculate_sha256(str(random_uuid))

def get_unique_id() -> str:
    return str(ULID())

def get_current_time() -> int:
    return int(time.time() * 1000) + datetime.now().microsecond // 1000

def error_return(code: int, message: str, data: T = None) -> GenericResponseModel:
    return GenericResponseModel(code=code, message=message, result=data)

def get_ts(date_string: str, format: str = "%Y/%m/%d %H:%M:%S") -> float:
    date_object = datetime.strptime(date_string, format)
    timestamp = date_object.timestamp()
    return timestamp

def get_time() -> datetime:
    tz = pytz.timezone('Asia/Shanghai')
    current_time = datetime.now(tz)
    return current_time

def success_return(data: T = None) -> GenericResponseModel:
    return GenericResponseModel(result=data)
