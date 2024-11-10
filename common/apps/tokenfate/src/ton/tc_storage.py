from pytonconnect.storage import IStorage
from datetime import datetime, timedelta
from typing import Optional, Tuple
from models.db_ops import FortunesDatabase
import redis.asyncio as redis
import json
import os
import logging
import asyncio
import os

REDIS_HOST = os.getenv("REDIS_HOST", "localhost")
REDIS_PORT = os.getenv("REDIS_PORT", 6379)
if not REDIS_HOST or not REDIS_PORT:
    raise Exception("REDIS_HOST and REDIS_PORT must be set")

client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT)

logging.info(f"Connected to Redis at {REDIS_HOST}:{REDIS_PORT}")


class TcStorage(IStorage):

    def __init__(self, chat_id: int):
        self.chat_id = chat_id

    def _get_key(self, key: str):
        return str(self.chat_id) + key

    async def set_item(self, key: str, value: str):
        print(f"Setting key {key} to {value}")
        await client.set(name=self._get_key(key), value=value)

    async def get_item(self, key: str, default_value: str = None):
        value = await client.get(name=self._get_key(key))
        print(f"Getting key {key}")
        return value.decode() if value else default_value

    async def remove_item(self, key: str):
        print(f"Removing key {key}")
        await client.delete(self._get_key(key))


tasks = {}  # 用于跟踪任务的内存字典

class TransactionManager:
    def __init__(self):
        self.redis_client = client
        self.tasks = tasks

    async def start_transaction(self, transaction_code: str, coro: asyncio.coroutine, expiration_time=180):
        """启动一个新事务并将其元数据存储在 Redis 中。"""
        task = asyncio.create_task(coro)
        self.tasks[transaction_code] = task  # 在内存中跟踪任务

        # 在 Redis 中存储事务元数据
        await self.redis_client.set(transaction_code, json.dumps({"status": "pending", "message": "Transaction started"}), ex=expiration_time)  # 设置 TTL 为 180 秒

        logging.info(f"Transaction {transaction_code} started")
        return task

    async def get_transaction_status(self, transaction_code: str):
        """从 Redis 中检索事务的状态。"""
        status = await self.redis_client.get(transaction_code)
        if status:
            return json.loads(status)
        else:
            return None

    async def update_transaction_status(self, transaction_code: str, status: str, message: str, expiration_time=180):
        """在 Redis 中更新事务的状态。"""
        await self.redis_client.set(transaction_code, json.dumps({"status": status, "message": message}), ex=expiration_time)  # 设置 TTL 为 180 秒

    async def cancel_transaction(self, transaction_code: str) -> bool:
        """取消正在运行的事务。"""
        task = self.tasks.get(transaction_code)
        if task:
            task.cancel()
            del self.tasks[transaction_code]  # 从内存中删除任务
            await self.redis_client.delete(transaction_code)  # 从 Redis 中删除任务元数据
            logging.info(f"Transaction {transaction_code} cancelled")
            return True
        else:
            logging.warning(f"Transaction {transaction_code} not found")
            return False

class DailyFortune:
    def __init__(self):
        # Redis连接配置
        self.redis_client = client

        # SQLite3连接配置
        self.db = FortunesDatabase()

    def _generate_daily_key(self, chat_id: str, ticker: str) -> str:
        """
        生成每日用户-ticker的唯一键
        格式: fortune:{chat_id}:{ticker.lower()}
        """
        return f"fortune:{chat_id}:{ticker.lower()}"

    def _get_next_midnight(self) -> datetime:
        """
        获取下一个午夜时间
        """
        now = datetime.now()
        next_midnight = now.replace(hour=0, minute=0, second=0, microsecond=0) + timedelta(days=1)
        return next_midnight
    
    async def get_cached_lot(self, chat_id: str, ticker: str) -> dict:
        """
        检查特定聊天的ticker的今日签是否已缓存
        """
        key = self._generate_daily_key(chat_id, ticker)
        cached_lot = await self.redis_client.get(key)
        if cached_lot is not None:
            logging.info(f"Hit cache for chat_id {chat_id}, ticker {ticker}")
            return json.loads(cached_lot) # 从Redis中获取缓存的签并转换为dict
        return {}

    async def get_daily_lot(self, chat_id: str, ticker: str) -> dict:
        try:
            key = self._generate_daily_key(chat_id, ticker)
            
            # 如果没有缓存,则从sqlite3中随机抽签
            result_of_draw = self.db.randomly_choose_sign_by_weight()
            
            # 存储到Redis,设置到下一个午夜过期
            next_midnight = self._get_next_midnight()
            expire_seconds = (next_midnight - datetime.now()).total_seconds()
            await self.redis_client.set(key, json.dumps(result_of_draw), ex=int(expire_seconds))
            
            logging.info(f"Draw a new lot for chat_id {chat_id}, ticker {ticker}: {result_of_draw}")
            return result_of_draw
            
        except Exception as e:
            logging.error(f"Error in get_daily_fortune: {str(e)}")
            # 如果Redis出错,仍然返回一个固定的签,确保同一天同一个ticker得到相同结果
            result_of_draw = self.db.randomly_choose_sign_by_weight()
            return result_of_draw
        
    async def clear_fortune(self, chat_id: str, ticker: str) -> bool:
        """
        清除特定聊天的ticker的今日签
        用于测试或特殊情况下重置
        """
        try:
            key = self._generate_daily_key(chat_id, ticker)
            return await self.redis_client.delete(key) > 0
        except Exception as e:
            self.logging.error(f"Error in clear_fortune: {str(e)}")
            return False