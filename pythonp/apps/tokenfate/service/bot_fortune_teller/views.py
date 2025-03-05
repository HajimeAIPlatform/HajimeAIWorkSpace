import logging
import asyncio
from io import BytesIO
from os import getenv

from flask import Blueprint, jsonify, request
from md2tgmd import escape
from telegram import Update, BotCommand, InlineKeyboardButton, InlineKeyboardMarkup, WebAppInfo,InlineQueryResultsButton
from PIL import Image
from telegram.ext import ApplicationBuilder, DictPersistence,CommandHandler

from pythonp.apps.tokenfate.service.bot_fortune_teller.commands import set_bot_commands_handler
from pythonp.apps.tokenfate.service.dify.views import chat_blocking_fortune_teller


# 获取Telegram Bot Token
telegram_bot_fortune_teller_token = getenv('TELEGRAM_BOT_FORTUNE_TELLER_TOKEN')
if not telegram_bot_fortune_teller_token:
    logging.error("Telegram bot token is not set in the environment")
    raise ValueError("Telegram bot token is not set in the environment")

# 创建Telegram应用
persistence = DictPersistence()
telegram_app = ApplicationBuilder().token(telegram_bot_fortune_teller_token).persistence(
    persistence).build()

async def run_bot_fortune_teller():
    await set_bot_commands_handler(telegram_app)
    await telegram_app.initialize()
    await telegram_app.start()

TOKEN_FATE_TWITTER = getenv('TOKEN_FATE_TWITTER')
TOKEN_FATE_GROUP = getenv('TOKEN_FATE_GROUP')


# 创建Flask Blueprint
bot_fortune_teller = Blueprint('bot_fortune_teller', __name__)

async def start(update):
    await update.message.reply_text("Fortune Teller is an intelligent LLM with built-in ads system!")

@bot_fortune_teller.route('/webhook', methods=['POST'])
# @run_async
async def webhook():
    chat_id = None
    try:
        body = request.get_json()

        update = Update.de_json(body, telegram_app.bot)

        if update.edited_message:
            return 'OK'

        if update.message.text == '/start':
            await start(update)
            return "OK"

        if update.message.photo:
            file_id = update.message.photo[-1].file_id
            logging.info(f"Images file id is {file_id}")
            file = await telegram_app.bot.get_file(file_id)
            logging.info("Image file found")
            bytes_array = await file.download_as_bytearray()
            bytesIO = BytesIO(bytes_array)
            logging.info("Images file as bytes")
            image = Image.open(bytesIO)
            logging.info("Image opened")

            prompt = 'Describe the image'

            if update.message.caption:
                prompt = update.message.caption
            logging.info(f"Prompt is {prompt}")

            text = "test"

            return {
                "method": "sendMessage",
                "chat_id": chat_id,
                "text": escape(text),
                "parse_mode": "MarkdownV2"
            }
        else:
            chat_response = chat_blocking_fortune_teller({
                "query": update.message.text,
                "user": chat_id
            })

            # 将响应文本赋值给 text 变量
            text = chat_response

            # 回复用户消息，并附带键盘按钮
            await update.message.reply_text(escape(text), parse_mode="MarkdownV2")

        return {
            "method": "sendMessage",
            "chat_id": chat_id,
            "text": escape(text),
            "parse_mode": "MarkdownV2"
        }
    except Exception as error:
        logging.error(f"Error Occurred: {error}")
        return {
            "method":
                "sendMessage",
            "chat_id":
                chat_id,
            "text":
                'Sorry, I am not able to generate content for you right now. Please try again later.'
        }
