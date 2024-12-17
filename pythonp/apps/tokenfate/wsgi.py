'''
Description: 
Author: Devin
Date: 2024-11-21 14:18:32
'''
import os
from dotenv import load_dotenv
from gunicorn.app.base import BaseApplication
from uvicorn.workers import UvicornWorker
import multiprocessing
import logging

# 加载环境变量
load_dotenv()
from pythonp.apps.tokenfate.email_service.email_sender import EmailMonitor

smtp_config = {
    'smtp_host': os.getenv('SMTP_HOST'),
    'smtp_port': os.getenv('SMTP_PORT'),
    'smtp_user': os.getenv('SMTP_USER'),
    'smtp_pass': os.getenv('SMTP_PASS'),
    'email_from': os.getenv('EMAIL_FROM'),
    'email_to': os.getenv('EMAIL_TO')
}

# 创建并启动监控器
monitor = EmailMonitor(smtp_config)
monitor.start()

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
        'workers': int(workers),
        'worker_class': UvicornWorker,
        'lifespan': 'off',
        'preload_app': True,  # Preload application to improve worker stability
        'max_requests': 1000,  # Restart workers periodically to prevent memory leaks
        'max_requests_jitter': 50,  # Random variance to prevent worker restart thundering herd
    }

    # 假设你的 Flask 应用在 `wsgi.py` 中定义为 `app`
    from pythonp.apps.tokenfate.app import app
    logging.info(f"current pwd: {os.getcwd()}")
    logging.info(f"Working config: using {workers} workers")
    StandaloneApplication(app, options).run()

if __name__ == "__main__":
    run_gunicorn()
