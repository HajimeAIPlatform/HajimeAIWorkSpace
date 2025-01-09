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
from telegram import Update, BotCommand, InlineKeyboardButton, InlineKeyboardMarkup, WebAppInfo, \
    InlineQueryResultsButton, InputMediaPhoto, Message, CallbackQuery
from telegram.ext import ApplicationBuilder, DictPersistence, CommandHandler

from pythonp.apps.tokenfate.service.dify.views import chat_blocking, chat_streaming, chat_workflow, chat_decode, chat_tarot
# import pythonp.apps.tokenfate.service.ton.views as ton_module
from pythonp.apps.tokenfate.service.bot3.commands import set_commands
from pythonp.apps.tokenfate.static.static import get_images_path
from pythonp.apps.tokenfate.utils.debug_tools import get_user_friendly_error_info

# 获取Telegram Bot Token
telegram_bot3_token = getenv('TELEGRAM_BOT3_TOKEN')
if not telegram_bot3_token:
    logging.error("Telegram bot3 token is not set in the environment")
    raise ValueError("Telegram bot3 token is not set in the environment")

# 创建Telegram应用
persistence = DictPersistence()
telegram_app = ApplicationBuilder().token(telegram_bot3_token).persistence(
    persistence).build()

async def run_bot3():
    await set_commands(telegram_app.bot)
    await telegram_app.initialize()
    await telegram_app.start()

# 创建Flask Blueprint
bot3 = Blueprint('bot3', __name__)

async def start(update: Update):
    await update.message.reply_text("Please send me your Twitter handle to get your Fortune today!")

async def handle_streaming_chat(chat_id, query):
    data = {"query": query}
    generator = chat_streaming(data)

    try:
        # 先发送一条初始消息，并获取 message_id
        initial_message = await telegram_app.bot.send_message(
            chat_id=chat_id,
            text="Generating content, please wait...",
        )
        message_id = initial_message.message_id
        current_message_content = initial_message.text
        message = ""
        # 使用 message_id 更新消息内容
        for chunk in generator:
            if chunk == "" or chunk is None: continue
            message = message + chunk
            escaped_chunk = escape(message)
            logging.debug(f"Escaped chunk: {escaped_chunk}")

            # 仅在新内容与当前内容不同时更新消息
            if message != current_message_content:
                await telegram_app.bot.edit_message_text(
                    chat_id=chat_id,
                    message_id=message_id,
                    text=escaped_chunk,
                    parse_mode="MarkdownV2"
                )
                current_message_content = message

    except Exception as e:
        logging.error(f"Error occurred during streaming: {e}")
        await telegram_app.bot.send_message(
            chat_id=chat_id,
            text='Sorry, I am not able to generate content for you right now. Please try again later.'
        )

@bot3.route('/webhook', methods=['POST'])
async def webhook():
    chat_id = None
    try:
        body = request.get_json()
        update = Update.de_json(body, telegram_app.bot)
        if update.edited_message:
            return 'OK'
        if update.message:
            inputs = {
                "input": update.message.text,
            }
            answer = chat_tarot(inputs)
            await update.message.reply_photo(
                photo=answer["url"],
                caption=escape(answer["text"]),
                parse_mode="MarkdownV2",
            )
            return jsonify({'status': 'ok'}), 200
    
        
    except Exception as e:
        logging.error(f"Error parsing update: {e}")
        return {
            "method":
                "sendMessage",
            "chat_id":
                chat_id,
            "text":
                'Sorry, I am not able to generate content for you right now. Please try again later.'
        }