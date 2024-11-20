import os
import logging

from dotenv import load_dotenv
import uvicorn
from flask import Flask
from asgiref.wsgi import WsgiToAsgi
from pythonp.common.logging.logger import setup_logging
# import nest_asyncio

# nest_asyncio.apply()

load_dotenv()

from pythonp.apps.tokenfate.src import blueprint as api
from pythonp.apps.tokenfate.models import setup_db
from pythonp.apps.tokenfate.src.binance.transaction_queue import start_transaction_processor
from pythonp.apps.tokenfate.src.binance.schedule import start_schedule_thread


def create_app():
    backend_app = Flask(__name__)
    setup_db(backend_app)
    backend_app.register_blueprint(api)
    backend_app.config.DEBUG = False

    @backend_app.get('/')
    def hello_world():
        return 'Hello, World!'

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
