import pymongo

from blog.blog_schema import MaterialModel, EmailSubscribeModel
from db.models import BaseDocument

from pydantic import BaseModel, Field
from datetime import datetime
from typing import Optional, List, Set, Any

from pythonp.apps.hajime_blog.utils.common import get_unique_id, get_current_time, success_return, error_return


class Tag(BaseDocument):
    id:str=Field(default_factory=get_unique_id)
    tag:str=""
    create_at: int = Field(default_factory=get_current_time)


    class Settings:
        name = "tag"

    @classmethod
    async def upsert(cls,tag):
        o = await cls.find_one({"tag":tag})
        if o is None:
            doc = {
                "tag":tag
            }
            await  cls(**doc).create()


class Material(BaseDocument):
    id:str=Field(default_factory=get_unique_id)
    title:str=""
    content: str
    tags:List[str]=[]
    covers:List[Any]=[]
    usernames:List[str]=[]
    keywords:List[str]=[]
    category:str="A"
    publish_at:int=0
    is_published:int=0

    create_at: int = Field(default_factory=get_current_time)


    class Settings:
        name = "material"


    @classmethod
    async def upsert(cls,form:MaterialModel):
        items = form.tags
        for tag in items:
            await Tag.upsert(tag)

        if form.id != "":
            ret=await Material.find_one({"_id":form.id}).update_one({"$set":{
                "content": form.content,
                "tags": form.tags,
                "covers": form.covers,
                "usernames": form.usernames,
                "keywords": form.keywords,
                "title":form.title
            }})
            if ret.modified_count ==1:
                return  success_return()
            else:
                return error_return(1,"no change")


        else:
            if form.publish_at>0:
                ts = form.publish_at
            else:
                ts = get_current_time()
            data = {
                "content": form.content,
                "tags": form.tags,
                "covers": form.covers,
                "usernames": form.usernames,
                "keywords": form.keywords,
                "publish_at":ts,
                "title": form.title
            }
            o=await Material(**data).create()
            return success_return(MaterialModel(**o.model_dump()))



class EmailSubscribe(BaseDocument):
    id:str=Field(default_factory=get_unique_id)
    email:str=""
    normalized_email:str=""
    create_at: int = Field(default_factory=get_current_time)


    class Settings:
        name = "subcribe"


    @classmethod
    async def upsert(cls,form:EmailSubscribeModel):

        data = {
            "email": form.email,
            "normalized_email":form.email.lower()
        }
        o = await EmailSubscribe.find_one({"normalized_email":form.email.lower()})
        if o is None:
            o=await EmailSubscribe(**data).create()
            return success_return(EmailSubscribe(**o.model_dump()))
        else:
            return error_return(1,"already subscribed")



