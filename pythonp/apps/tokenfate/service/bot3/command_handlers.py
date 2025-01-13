from md2tgmd import escape
from pythonp.apps.tokenfate.service.dify.views import chat_tarot
from telegram import Update
from pythonp.apps.tokenfate.service.ton.tc_storage import TarotStorage
from pythonp.apps.tokenfate.service.bot3.checks_before_handler import drawn_card, check_amount, check_before_tarot
from datetime import datetime
import logging
from flask import jsonify
from telegram import InlineKeyboardButton, InlineKeyboardMarkup

from telegram.ext import (
    ContextTypes,
)

# 定义每个命令的处理逻辑
async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    text="Greetings, seeker. I am Tarot Sakura, your guide into the arcane world of tarot cards. Please tell me which area of your life you wish to explore or ask a specific question so that I may draw cards accordingly. Will it be about Love? Career? Personal growth? Or perhaps, a different question altogether? The cards are ready to reveal their secrets."
    await update.message.reply_text(text=escape(text), parse_mode="MarkdownV2")

async def tarot(update: Update, context: ContextTypes.DEFAULT_TYPE):
    if not await check_before_tarot(update, context):
        return jsonify({'status': 'ok'})
    tarot_storage = TarotStorage()
    if not context.args:
        # 如果没有参数，提示用户 /tarot 的使用方法
        text="Usage: /tarot question"
        await update.message.reply_text(
            text=escape(text),
            parse_mode="MarkdownV2",
        )
    else:
        # 如果有参数，执行 chat_tarot 处理逻辑
        user_question = " ".join(context.args)  # 合并所有参数为完整问题
        inputs = {
            "input": user_question,
        }
        answer = chat_tarot(inputs)
        await drawn_card(update, context)
        await update.message.reply_photo(
            photo=answer["url"],
            caption=escape(answer["text"]),
            parse_mode="MarkdownV2",
        )
        today_history = {
            "Question": user_question,
            "Answer": answer["text"],
            # "Cards": answer["cards"],
            "Time": datetime.now().strftime("%Y-%m-%d %H:%M:%S"), # 当前时间 
        }
        await tarot_storage.store_today_draw(update.message.chat_id, today_history)

async def history(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    tarot_storage = TarotStorage()
    if not await tarot_storage.is_today_drawn(update.message.chat_id):
        text="Today's tarot has not been drawn yet. Please use /tarot to draw today's tarot first."
        await update.message.reply_text(
            text=escape(text),
            parse_mode="MarkdownV2",
        )
    else:
        today_history = await tarot_storage.get_today_draw(update.message.chat_id)
        logging.info(today_history)
        text = f"**Today's tarot history**\n\nQuestion:\n {today_history['Question']}\nAnswer:\n {today_history['Answer']}\nTime: {today_history['Time']}"
        await update.message.reply_text(text=escape(text), parse_mode="MarkdownV2")

async def community(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    keyboard = [[InlineKeyboardButton(
                        text="HajimeAI 🤖 Official Community",
                        url="t.me/HajimeAI"
                    )]]
    reply_markup = InlineKeyboardMarkup(keyboard)
    await update.message.reply_text("Join our community!", reply_markup=reply_markup)

async def amount(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await check_amount(update, context)
