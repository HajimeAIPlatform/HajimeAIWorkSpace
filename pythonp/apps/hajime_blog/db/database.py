from beanie import init_beanie
from pythonp.apps.hajime_blog.config.connection import db
from pythonp.apps.hajime_blog.db.models import *
# EvidenceData, UserDeposit, LastCheckTime, User

from pythonp.apps.hajime_blog.blog.blog_model import *

async def init_db():
    await init_beanie(database=db, document_models=[EvidenceData,UserDeposit,LastCheckTime,
                                                    User,UserStat,UserAssetList,UserAsset,UserWithdraw,BizUser,MinerOrder,
                                                    UserNewsSubscribe,OpLog,
                                                    Material,
                                                    EmailSubscribe,
                                                    Tag
                                                    ])
    counter_collection = await db.counter.find_one({"type": "user"})
    if counter_collection is None:
        await db.create_collection("counter")
        await db.counter.insert_one({"type": "user", "seq": 0})



