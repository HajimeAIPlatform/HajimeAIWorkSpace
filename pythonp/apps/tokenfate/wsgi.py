import os
from dotenv import load_dotenv
from gunicorn.app.base import BaseApplication
from uvicorn.workers import UvicornWorker
import multiprocessing
import logging

# 加载环境变量
load_dotenv()

class StandaloneApplication(BaseApplication):
    def __init__(self, app, options=None):
        self.options = options or {}
        self.application = app
        super().__init__()

    def load_config(self):
        config = {key: value for key, value in self.options.items()
                  if key in self.cfg.settings and value is not None}
        for key, value in config.items():
            self.cfg.set(key.lower(), value)

    def load(self):
        return self.application

def run_gunicorn():
    # 从环境变量获取配置
    workers_env = os.getenv('WORKERS')
    host = os.getenv('HOST', '0.0.0.0')
    port = os.getenv('PORT', '5000')

    if workers_env:
        workers = int(workers_env)
    else:
        workers = multiprocessing.cpu_count() * 2 + 1
  
    options = {
        'bind': f'{host}:{port}',
        'workers': int(workers * 2 + 1),
        'worker_class': UvicornWorker,
        'lifespan': 'off',  # 禁用 lifespan
    }

    # 假设你的 Flask 应用在 `wsgi.py` 中定义为 `app`
    from app import app
    logging.info(f"Working config: using {workers} workers")
    StandaloneApplication(app, options).run()

if __name__ == "__main__":
    run_gunicorn()
