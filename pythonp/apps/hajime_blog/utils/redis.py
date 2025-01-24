import os
from redis import asyncio as aioredis
from dotenv import load_dotenv

load_dotenv()

redis_client = None

class RedisClient:
    def __init__(self):
        self.url = os.getenv("REDIS_URL", "redis://localhost")
        self.db = int(os.getenv("REDIS_DB", 0))
        self.decode_responses = os.getenv("REDIS_DECODE_RESPONSES", "true").lower() == "true"
        self.password = os.getenv("REDIS_PASSWORD")
        self._client = None

    async def get_client(self):
        if not self._client:
            self._client = await aioredis.from_url(self.url, db=self.db,password=self.password, decode_responses=self.decode_responses)
        return self._client

    async def close(self):
        if self._client:
            await self._client.close()
            self._client = None


async def get_redis():
    global redis_client
    if not redis_client:
        redis_client = RedisClient()
    
    return await redis_client.get_client() 
