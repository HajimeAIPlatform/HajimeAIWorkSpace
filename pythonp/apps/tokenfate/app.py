import os
import logging

from dotenv import load_dotenv
import uvicorn
from flask import Flask
from asgiref.wsgi import WsgiToAsgi
from pythonp.common.logging.logger import setup_logging
import asyncio
# import nest_asyncio

# nest_asyncio.apply()

load_dotenv()

from pythonp.apps.tokenfate.service import blueprint as api
from pythonp.apps.tokenfate.models import setup_db
from pythonp.apps.tokenfate.service.binance.transaction_queue import start_transaction_processor
from pythonp.apps.tokenfate.service.binance.schedule import start_schedule_thread
from pythonp.apps.tokenfate.service.bot.views import run_bot
from pythonp.apps.tokenfate.service.bot3.views import run_bot3



def create_app():
    backend_app = Flask(__name__)
    setup_db(backend_app)
    backend_app.register_blueprint(api)
    backend_app.config.DEBUG = False

    @backend_app.get('/')
    def hello_world():
        return 'Hello, World!'

    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    async def run_all_bots():
        await asyncio.gather(run_bot(), run_bot3())
    
    # Run the bot asynchronously
    loop.run_until_complete(run_all_bots())

    start_transaction_processor(backend_app)
    start_schedule_thread()

    return WsgiToAsgi(backend_app)

# 从环境变量中获取端口号，如果没有设置则使用默认的5000
port = int(os.getenv("PORT", 5000))
env = os.getenv("ENV", "dev")

log_dir = os.getenv("LOG_DIR", "./logs")
setup_logging(log_dir=log_dir)
if env == "prod":
    logging.disable(logging.DEBUG)
    logging.info("Running in production mode")
else:
    logging.disable(logging.DEBUG)
    logging.info("Running in development mode")

app = create_app()


# if __name__ == "__main__":
#     # 运行 Flask 应用
#     logging.info(f"Starting backend server on port {port}...")
#     uvicorn.run(app,
#                 host="0.0.0.0",
#                 port=port,
#                 log_level="info",
#                 access_log=False)
