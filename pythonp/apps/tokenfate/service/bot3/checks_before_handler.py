from telegram import (
    Update, 
)
from telegram.ext import (
    ContextTypes,
)
from pythonp.apps.tokenfate.models.transaction import TarotUser
from md2tgmd import escape

async def check_status(update: Update, context: ContextTypes.DEFAULT_TYPE):
    # check all status
    await check_first_visit(update, context)
    await check_sign_in(update, context)

async def check_first_visit(update, context: ContextTypes.DEFAULT_TYPE):
    # check first visit
    chat_id = update.effective_chat.id
    text="Welcome to Tarot Bot, +100 points for your first visit!"
    result = TarotUser.create_user(str(chat_id))
    if result.get("status") == "success":
        await context.bot.send_message(
            chat_id=chat_id,
            text=escape(text),
            parse_mode="MarkdownV2",
        )

async def check_sign_in(update: Update, context: ContextTypes.DEFAULT_TYPE):
    # check sign in
    chat_id = update.effective_chat.id
    result = TarotUser.sign_in(str(chat_id))
    text="+20 points for signing in successfully!"
    if result.get("status") == "success":
        await context.bot.send_message(
            chat_id=chat_id,
            text=escape(text),
            parse_mode="MarkdownV2",
        )
    
