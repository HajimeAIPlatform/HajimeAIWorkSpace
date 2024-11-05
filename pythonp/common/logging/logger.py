import os
import logging
from datetime import datetime, timezone
from logging.handlers import TimedRotatingFileHandler


class CustomTimedRotatingFileHandler(TimedRotatingFileHandler):

    def __init__(self, log_dir, *args, **kwargs):
        self.log_dir = log_dir
        super().__init__(*args, **kwargs)

    def doRollover(self):
        current_date = datetime.now(timezone.utc).strftime("%Y-%m-%d")
        log_path = os.path.join(self.log_dir, current_date)
        if not os.path.exists(log_path):
            os.makedirs(log_path)

        self.baseFilename = os.path.join(log_path, 'bot_log.log')
        super().doRollover()


def setup_logging(log_dir):
    current_date = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    log_path = os.path.join(log_dir, current_date)
    if not os.path.exists(log_path):
        os.makedirs(log_path)

    log_file = os.path.join(log_path, 'bot_log.log')
    handler = CustomTimedRotatingFileHandler(log_dir,
                                             log_file,
                                             when="midnight",
                                             interval=1)
    handler.suffix = "%Y-%m-%d"

    formatter = logging.Formatter(
        '%(asctime)s - %(pathname)s - %(levelname)s - %(funcName)s - line: %(lineno)d - %(message)s'
    )
    handler.setFormatter(formatter)
    handler.setLevel(logging.DEBUG)

    logger = logging.getLogger('')
    logger.setLevel(logging.DEBUG)
    logger.addHandler(handler)

    console = logging.StreamHandler()
    console.setLevel(logging.INFO)
    console.setFormatter(formatter)
    logger.addHandler(console)
