from fastapi import FastAPI
from blog.blog_router import router as blog_router
from utils.common import get_redis

app = FastAPI()

# 挂载博客模块的路由
app.include_router(blog_router, prefix="/blog", tags=["Blog"])

# 启动事件
@app.on_event("startup")
async def startup():
    app.state.redis = await get_redis()
    print("Redis connected!")

# 关闭事件
@app.on_event("shutdown")
async def shutdown():
    redis = app.state.redis
    if redis:
        await redis.close()
        print("Redis connection closed.")

# 测试路由
@app.get("/")
def read_root():
    return {"message": "Hajime Blog API is running!"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="127.0.0.1", port=8000)
