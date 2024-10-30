import asyncio
from asyncio import sleep

import pymysql
from fastapi import FastAPI, Depends, HTTPException, BackgroundTasks

from apscheduler.schedulers.asyncio import AsyncIOScheduler
from fastapi.exceptions import RequestValidationError
from starlette.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware

from app.database import init_db
from app.db.models import EvidenceData, UserDeposit
from app.db.schemas import EvidenceDataInputModel, GenericResponseModel, SignLoginModel
from app.routers import api_router, admin_router, admin_auth_router
from app.service.solana_service import SolanaService
from app.blog  import blog_router

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

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

@app.on_event("startup")
async def app_start():
    scheduler.add_job(execute_periodic_function, 'interval', seconds=30)
    scheduler.start()
    await init_db()


@app.on_event("shutdown")
async def shutdown_event():
    scheduler.shutdown()


async def save_hash_to_blockchain(task_id: str, data_hash: str,node_id:str):
    transaction_hash = await SolanaService.call_evidence_contract(data_hash)
    print("task_id,block_hash:",task_id, transaction_hash)
    if transaction_hash != "":
        await EvidenceData.update_task_hash(task_id, transaction_hash)
        await SolanaService.callback(node_id,transaction_hash)



async def background_task(task_id, data_hash,node_id):
    asyncio.create_task(save_hash_to_blockchain(task_id, data_hash,node_id))



@app.get("/ping")
def ping():
    return GenericResponseModel()

