from telegram.ext import ContextTypes
import logging

# 添加错误处理器到你的 telegram_app
async def error_handler(update: object, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Log the error and send a message to inform the developer."""
    logging.error(msg="Exception while handling an update:", exc_info=context.error)