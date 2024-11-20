import os
from dotenv import load_dotenv
from gunicorn.app.base import BaseApplication
from uvicorn.workers import UvicornWorker

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
    workers = os.getenv('WORKERS', '4')
    host = os.getenv('HOST', '0.0.0.0')
    port = os.getenv('PORT', '8000')

    options = {
        'bind': f'{host}:{port}',
        'workers': int(workers),
        'worker_class': UvicornWorker,
    }

    # 假设你的 Flask 应用在 `wsgi.py` 中定义为 `app`
    from wsgi import app
    StandaloneApplication(app, options).run()

if __name__ == "__main__":
    run_gunicorn()
