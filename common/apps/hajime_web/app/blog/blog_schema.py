import pymongo

from app.db.models import BaseDocument

from pydantic import BaseModel, Field
from datetime import datetime
from typing import Optional, List, Set, Any

from app.utils.common import get_unique_id, get_current_time
class EmailSubscribeModel(BaseModel):
    email: str

class MaterialModel(BaseModel):
    id: str = ""
    content: str
    title: str = ""
    tags: List[str] = []
    keywords: List[str] = []
    covers: List[Any] = []
    usernames: list[str] = []
    publish_at: int = 0


class AdminMaterialListRequestModel(BaseModel):
    page: int = 0
    pagesize: int = 100
    tag: str = ""



class AdminMaterialListItemModel(BaseModel):
    id: str
    title: str = ""
    content: str
    tags: List[str] = []
    covers: List[Any] = []
    usernames: list[str] = []
    keywords: list[str] = []
    publish_at: int = 0
    is_published:int=0

class AdminTagListItemModel(BaseModel):
    id:str
    tag:str
    create_at:int


class PublishTypeModel(BaseModel):
    type:str

