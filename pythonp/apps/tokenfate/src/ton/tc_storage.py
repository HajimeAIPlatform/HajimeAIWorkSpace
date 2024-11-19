import json
import os
import logging
import asyncio
import os
from datetime import datetime, timedelta
from typing import Optional, Tuple, Set

from pytonconnect.storage import IStorage
import redis.asyncio as redis

from  pythonp.apps.tokenfate.models.db_ops import FortunesDatabase

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

    def _generate_key(self, chat_id: str, ticker: str, prefix: str, lang: str = None) -> str:
        if lang:
            return f"{prefix}:{chat_id}:{ticker.lower()}:{lang}"
        return f"{prefix}:{chat_id}:{ticker.lower()}"

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
        key = self._generate_key(chat_id, ticker, "fortune")
        cached_lot = await self.redis_client.get(key)
        if cached_lot is not None:
            logging.info(f"Hit cache for chat_id {chat_id}, ticker {ticker}")
            return json.loads(cached_lot) # 从Redis中获取缓存的签并转换为dict
        return None
    
    async def get_cached_decode(self, chat_id: str, ticker: str, lang: str) -> str:
        key = self._generate_key(chat_id, ticker, "decode", lang)
        cached_decode = await self.redis_client.get(key)
        if cached_decode is not None:
            logging.info(f"Hit decode cache for chat_id {chat_id}, ticker {ticker}, lang {lang}")
            return cached_decode.decode('utf-8')  # 将字节串解码为字符串
        return None
    
    async def set_decode_cache(self, chat_id: str, ticker: str, lang: str, decode: str) -> None:
        key = self._generate_key(chat_id, ticker, "decode", lang)
        next_midnight = self._get_next_midnight()
        expire_seconds = (next_midnight - datetime.now()).total_seconds()
        await self.redis_client.set(key, decode, ex=int(expire_seconds))
        logging.info(f"Cached decode for chat_id {chat_id}, ticker {ticker}, lang {lang}: {decode}")
        # 设置解签状态
        await self.set_decode_status(chat_id, ticker)

    async def get_decode_status(self, chat_id: str, ticker: str) -> bool:
        key = self._generate_key(chat_id, ticker, "decode_status")
        status = await self.redis_client.get(key)
        if status is not None:
            logging.info(f"Decode status already recorded for chat_id {chat_id}, ticker {ticker}")
            return True
        return False

    async def set_decode_status(self, chat_id: str, ticker: str) -> None:
        key = self._generate_key(chat_id, ticker, "decode_status")
        next_midnight = self._get_next_midnight()
        expire_seconds = (next_midnight - datetime.now()).total_seconds()
        await self.redis_client.set(key, "1", ex=int(expire_seconds))
        logging.info(f"Recorded decode status for chat_id {chat_id}, ticker {ticker}")

    async def get_daily_lot(self, chat_id: str, ticker: str) -> dict:
        try:
            key = self._generate_key(chat_id, ticker, "fortune")
            
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
            key = self._generate_key(chat_id, ticker, "fortune")
            return await self.redis_client.delete(key) > 0
        except Exception as e:
            self.logging.error(f"Error in clear_fortune: {str(e)}")
            return False
        
class UserActivityTracker:
    def __init__(self):
        self.redis_client = client
        
    def _get_wallet_set_key(self, user_id: str) -> str:
        """获取用户已连接钱包集合的Redis键"""
        return f"user:{user_id}:connected_wallets"
        
    def _get_checkin_key(self, user_id: str) -> str:
        """获取用户打卡记录的Redis键"""
        return f"user:{user_id}:last_checkin"
    
    async def _has_checked_in_today(self, user_id: str) -> bool:
        """检查用户今天是否已经打卡"""
        checkin_key = self._get_checkin_key(user_id)
        return await self.redis_client.exists(checkin_key) 

    def _get_next_midnight(self) -> datetime:
        """获取下一个午夜时间"""
        now = datetime.now()
        next_midnight = now.replace(hour=0, minute=0, second=0, microsecond=0) + timedelta(days=1)
        return next_midnight

    async def connect_wallet(self, user_id: str, wallet_name: str) -> bool:
        """
        处理用户连接钱包的逻辑
        
        Args:
            user_id: 用户ID
            wallet_name: 钱包名称
            
        Returns:
            bool: 是否是首次连接这种钱包
        """
        try:
            wallet_set_key = self._get_wallet_set_key(user_id)
            
            # 检查是否首次连接这种钱包
            is_new_wallet = await self.redis_client.sadd(wallet_set_key, wallet_name)
            
            if is_new_wallet:
                logging.info(f"User {user_id} connected new wallet type: {wallet_name}")
            else:
                logging.info(f"User {user_id} reconnected existing wallet type: {wallet_name}")
            
            return bool(is_new_wallet)
            
        except Exception as e:
            logging.error(f"Error in connect_wallet: {str(e)}")
            raise
    
    async def disconnect_wallet(self, user_id: str, wallet_name: str) -> bool:
        """撤回用户连接的钱包"""
        try:
            wallet_set_key = self._get_wallet_set_key(user_id)
            removed = await self.redis_client.srem(wallet_set_key, wallet_name)
            
            if removed:
                logging.info(f"User {user_id} disconnected wallet type: {wallet_name}")
            else:
                logging.info(f"Wallet type {wallet_name} was not connected for user {user_id}")
            
            return bool(removed)
            
        except Exception as e:
            logging.error(f"Error in disconnect_wallet: {str(e)}")
            raise

    async def daily_checkin(self, user_id: str) -> bool:
        """处理用户每日打卡"""
        try:
            if await self._has_checked_in_today(user_id):
                return False
            
            # 未打卡，记录打卡时间
            checkin_key = self._get_checkin_key(user_id)
            next_midnight = self._get_next_midnight()
            expire_seconds = int((next_midnight - datetime.now()).total_seconds())
            
            # 设置打卡记录，到午夜过期
            await self.redis_client.set(
                checkin_key,
                str(datetime.now().timestamp()),
                ex=expire_seconds
            )
            
            logging.info(f"User {user_id} checked in successfully")
            return True
            
        except Exception as e:
            logging.error(f"Error in daily_checkin: {str(e)}")
            raise
            
    async def get_connected_wallets(self, user_id: str) -> Set[str]:
        """获取用户已连接的钱包列表"""
        try:
            wallet_set_key = self._get_wallet_set_key(user_id)
            wallets = await self.redis_client.smembers(wallet_set_key)
            return {w.decode() for w in wallets}
        except Exception as e:
            logging.error(f"Error in get_connected_wallets: {str(e)}")
            return set()
            
    async def is_checked_in_today(self, user_id: str) -> bool:
        """检查用户今天是否已经打卡"""
        try:
            return await self._has_checked_in_today(user_id)
        except Exception as e:
            logging.error(f"Error in is_checked_in_today: {str(e)}")
            return False