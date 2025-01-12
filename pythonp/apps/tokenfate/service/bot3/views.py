import logging
import asyncio
from io import BytesIO
from os import getenv
import urllib.parse
import time
import traceback
import json
import os
import re
from typing import List, Dict, Union
import sys
from flask import Blueprint, jsonify, request
from md2tgmd import escape
from telegram import (
    Update, 
)
from telegram.ext import (
    ApplicationBuilder, 
    DictPersistence, 
    ContextTypes,
    CommandHandler, 
    MessageHandler, 
    filters
)

# import pythonp.apps.tokenfate.service.ton.views as ton_module
from pythonp.apps.tokenfate.service.bot3.commands import set_commands
import pythonp.apps.tokenfate.service.bot3.command_handlers as command_handlers
from pythonp.apps.tokenfate.service.bot3.checks_before_handler import check_status
from pythonp.apps.tokenfate.service.bot3.message_handlers import reply_chat_tarot

# 获取Telegram Bot Token
telegram_bot3_token = getenv('TELEGRAM_BOT3_TOKEN')
if not telegram_bot3_token:
    logging.error("Telegram bot3 token is not set in the environment")
    raise ValueError("Telegram bot3 token is not set in the environment")

# 创建Telegram应用
persistence = DictPersistence()
telegram_app = ApplicationBuilder().token(telegram_bot3_token).persistence(persistence).build()

async def run_bot3():
    try:
        logging.info("Setting commands for the bot tarot...")
        await set_commands(telegram_app.bot)
        
        logging.info("Initializing the bot tarot...")
        await telegram_app.initialize()
        
        logging.info("Starting the bot tarot...")
        await telegram_app.start()

    except Exception as e:
        logging.error(f"An error occurred: {e}")

# 创建Flask Blueprint
bot3 = Blueprint('bot3', __name__)

# 注册消息处理器
def register_handlers(telegram_app):
    # 状态检查处理器 (较高优先级组)
    telegram_app.add_handler(
        MessageHandler(filters.ALL, check_status), 
        group=1
    )

    # 命令处理器
    telegram_app.add_handler(CommandHandler('start', command_handlers.start))
    telegram_app.add_handler(CommandHandler('history', command_handlers.history)) 
    telegram_app.add_handler(CommandHandler('community', command_handlers.community))
    telegram_app.add_handler(CommandHandler('integral', command_handlers.integral))
    
    # 文本消息处理器 - 需要放在命令处理器之后
    telegram_app.add_handler(MessageHandler(
        filters.TEXT & ~filters.COMMAND,
        reply_chat_tarot
    ))


# 初始化时注册所有处理器
register_handlers(telegram_app)

@bot3.route('/webhook', methods=['POST'])
async def webhook():
    chat_id = None
    try:
        body = request.get_json()
        update = Update.de_json(body, telegram_app.bot)

        # 统一使用 process_update 处理所有更新
        await telegram_app.process_update(update)

        # 确保总是返回一个响应
        return jsonify({'status': 'ok'})
    
    except Exception as e:
        logging.error(f"Error processing update: {e}")
        return {
            "method":
                "sendMessage",
            "chat_id":
                chat_id,
            "text":
                'Sorry, I am not able to generate content for you right now. Please try again later.'
        }