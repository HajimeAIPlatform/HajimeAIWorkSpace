import logging
import asyncio
from io import BytesIO
from os import getenv
import threading

from flask import Flask, Blueprint, jsonify, request
from md2tgmd import escape
from telegram import Update, BotCommand, InlineKeyboardButton, InlineKeyboardMarkup, InputFile
from PIL import Image
from telegram.ext import ApplicationBuilder, CommandHandler, MessageHandler, CallbackQueryHandler
import qrcode
from pytoniq_core import Address

from pythonp.apps.fortune_teller.service.bot.commands import set_bot_commands_handler
from pythonp.apps.fortune_teller.service.bot.wallet_menu_callback import set_handlers
from pythonp.apps.fortune_teller.service.dify.views import chat_blocking, chat_streaming
from pythonp.apps.fortune_teller.service.binance.views import handle_binance_command
import pythonp.apps.fortune_teller.service.ton.views as ton_module
from pythonp.apps.fortune_teller.utils.utils import run_async

# 获取Telegram Bot Token
telegram_bot_token = getenv('TELEGRAM_BOT_FORTUNE_TELLER_TOKEN')
if not telegram_bot_token:
    logging.error("Telegram bot token is not set in the environment")
    raise ValueError("Telegram bot token is not set in the environment")

# 创建Telegram应用
telegram_app = ApplicationBuilder().token(telegram_bot_token).build()


# 设置Bot命令
# asyncio.run(set_bot_commands_handler(telegram_app))
async def handle_streaming_chat(chat_id, query):
    data = {"query": query}
    generator = chat_streaming(data)

    try:
        for chunk in generator:
            await telegram_app.bot.send_message(chat_id=chat_id,
                                                text=escape(chunk),
                                                parse_mode="MarkdownV2")
    except Exception as error:
        logging.error(f"Error occurred during streaming: {error}")
        await telegram_app.bot.send_message(
            chat_id=chat_id,
            text=
            'Sorry, I am not able to generate content for you right now. Please try again later.'
        )


async def ton_command_handle(update):
    try:
        ton_response = await ton_module.handle_ton_command(
            telegram_app, update)
        if ton_response:
            return ton_response
        return None
    except Exception as error:
        logging.error(f"Error occurred during streaming: {error}")
        return {
            "text":
            'Sorry, I am not able to generate content for you right now. Please try again later.'
        }


async def start(update, context):
    await context.bot.send_message(chat_id=update.effective_chat.id,
                                   text="Welcome!")


async def handle_message(update, context):
    chat_id = update.message.chat_id
    text = update.message.text

    binance_response = handle_binance_command(text)
    if binance_response:
        await context.bot.send_message(chat_id=chat_id, text=binance_response)
        return

    if update.message.photo:
        logging.info('Generating images')
        file_id = update.message.photo[-1].file_id
        logging.info(f"Images file id is {file_id}")
        file = await context.bot.get_file(file_id)
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

        await context.bot.send_message(chat_id=chat_id,
                                       text=escape(text),
                                       parse_mode="MarkdownV2")
        return

    chat_response = chat_blocking({"query": text})
    logging.info(f"Chat response is {chat_response}")

    await context.bot.send_message(chat_id=chat_id,
                                   text=escape(chat_response),
                                   parse_mode="MarkdownV2")


async def handle_callback_query(update, context):
    await set_handlers(update, telegram_app)


async def start_telegram_bot():
    # 注册处理程序
    await set_bot_commands_handler(telegram_app)
    telegram_app.add_handler(CommandHandler("start", start))
    telegram_app.add_handler(CallbackQueryHandler(handle_callback_query))
    telegram_app.run_polling()


# 启动轮询
asyncio.run(start_telegram_bot())
