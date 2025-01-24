import asyncio
from asyncio import sleep
from fastapi import FastAPI, Depends, HTTPException, BackgroundTasks

from apscheduler.schedulers.asyncio import AsyncIOScheduler
from fastapi.exceptions import RequestValidationError
from starlette.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware


from pythonp.apps.hajime_blog.utils.redis import get_redis  

from pythonp.apps.hajime_blog.db.database import init_db
from pythonp.apps.hajime_blog.db.models import EvidenceData, UserDeposit
from pythonp.apps.hajime_blog.db.schemas import EvidenceDataInputModel, GenericResponseModel, SignLoginModel
from pythonp.apps.hajime_blog.routers import api_router, admin_router, admin_auth_router
from pythonp.apps.hajime_blog.blog import blog_router
from pythonp.apps.hajime_blog.service.solana_service import SolanaService


app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 挂载博客模块的路由
app.include_router(api_router.router)
app.include_router(admin_router.router)
app.include_router(admin_auth_router.router)
app.include_router(blog_router.router)

# Custome Excpeption Handler
@app.exception_handler(HTTPException)
async def http_exception_handler(request, exc):
    return JSONResponse(
        status_code=200,
        content={"code": exc.status_code, "message": exc.detail}
    )



@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request, exc):
    return JSONResponse(
        status_code=200,
        content={"code": 400, "message": "Validation error", "details": exc.errors()}
    )


async def execute_periodic_function():
    """
    TODO: Check failed task
    :return:
    """
    pass

scheduler = AsyncIOScheduler()

# 启动事件
@app.on_event("startup")
async def startup():
    scheduler.add_job(execute_periodic_function, 'interval', seconds=30)
    scheduler.start()

    app.state.redis = await get_redis()  # 获取 Redis 客户端
    print("Redis connected!")

# 关闭事件
@app.on_event("shutdown")
async def shutdown():
    scheduler.shutdown()
    redis = app.state.redis
    if redis:
        await redis.close()  # 关闭 Redis 连接
        print("Redis connection closed.")

async def save_hash_to_blockchain(task_id: str, data_hash: str,node_id:str):
    transaction_hash = await SolanaService.call_evidence_contract(data_hash)
    print("task_id,block_hash:",task_id, transaction_hash)
    if transaction_hash != "":
        await EvidenceData.update_task_hash(task_id, transaction_hash)
        await SolanaService.callback(node_id,transaction_hash)


async def background_task(task_id, data_hash,node_id):
    asyncio.create_task(save_hash_to_blockchain(task_id, data_hash,node_id))

# 测试路由
@app.get("/")
def read_root():
    return {"message": "Hajime Blog API is running!"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="127.0.0.1", port=8000)
