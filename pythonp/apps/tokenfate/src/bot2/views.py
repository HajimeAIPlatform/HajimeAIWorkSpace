import logging
import asyncio
from io import BytesIO
from os import getenv

from flask import Blueprint, jsonify, request
from md2tgmd import escape
from telegram import Update, BotCommand, InlineKeyboardButton, InlineKeyboardMarkup, WebAppInfo,InlineQueryResultsButton
from PIL import Image
from telegram.ext import ApplicationBuilder, DictPersistence,CommandHandler

from pythonp.apps.tokenfate.src.bot2.commands import set_bot_commands_handler
from pythonp.apps.tokenfate.src.dify.views import chat_blocking_key_2


# 获取Telegram Bot Token
telegram_bot2_token = getenv('TELEGRAM_BOT2_TOKEN')
if not telegram_bot2_token:
    logging.error("Telegram bot token is not set in the environment")
    raise ValueError("Telegram bot token is not set in the environment")

# 创建Telegram应用
persistence = DictPersistence()
telegram_app = ApplicationBuilder().token(telegram_bot2_token).persistence(
    persistence).build()

async def run_bot():
    await set_bot_commands_handler(telegram_app)
    await telegram_app.initialize()
    await telegram_app.start()

TOKEN_FATE_TWITTER = getenv('TOKEN_FATE_TWITTER')
TOKEN_FATE_GROUP = getenv('TOKEN_FATE_GROUP')

asyncio.run(run_bot())

# 创建Flask Blueprint
bot2 = Blueprint('bot2', __name__)

async def start(update):
    await update.message.reply_text("Please send me your Twitter handle to get your Fortune today!")


@bot2.route('/webhook', methods=['POST'])
# @run_async
async def webhook():
    chat_id = None
    try:
        body = request.get_json()

        update = Update.de_json(body, telegram_app.bot)

        print(update, 'update')

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
            chat_response = chat_blocking_key_2({
                "query": update.message.text,
                "user": chat_id
            })

            # 将响应文本赋值给 text 变量
            text = chat_response

            # 创建两个按钮
            button1 = InlineKeyboardButton(text="TokenFate Twitter", url=TOKEN_FATE_TWITTER)
            button2 = InlineKeyboardButton(text="Join TokenFate Group", url=TOKEN_FATE_GROUP)

            # 创建键盘布局
            keyboard = InlineKeyboardMarkup([[button1], [button2]])

            more_read = f"\nFor further analysis, please follow [TokenFate Twitter]({TOKEN_FATE_TWITTER}) or join the [TokenFate group]({TOKEN_FATE_GROUP}).\n"
            text = chat_response + '\n' + more_read

            # 回复用户消息，并附带键盘按钮
            await update.message.reply_text(escape(text), parse_mode="MarkdownV2", reply_markup=keyboard)

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
