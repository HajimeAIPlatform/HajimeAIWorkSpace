from beanie import init_beanie
from app.connection import db
from app.db.models import *
# EvidenceData, UserDeposit, LastCheckTime, User

from app.blog.blog_model import *
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





