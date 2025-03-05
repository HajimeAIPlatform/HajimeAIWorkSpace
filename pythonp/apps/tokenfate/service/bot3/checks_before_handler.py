from telegram import (
    Update, 
)
from telegram.ext import (
    ContextTypes,
)
from pythonp.apps.tokenfate.models.transaction import TarotUser
from md2tgmd import escape
from pythonp.apps.tokenfate.service.ton.tc_storage import TarotStorage

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

async def drawn_card(update: Update, context: ContextTypes.DEFAULT_TYPE):
    # check drawn card
    chat_id = update.effective_chat.id
    result = TarotUser.draw_card(str(chat_id))
    text="-20 points for drawing a card!"
    if result.get("status") == "success":
        await context.bot.send_message(
            chat_id=chat_id,
            text=escape(text),
            parse_mode="MarkdownV2",
        )

async def check_amount(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = update.effective_chat.id
    result = TarotUser.get_user_info(str(chat_id))
    if result.get("status") == "success":
        text = f"Your current balance is {result.get('amount')} points."
        await context.bot.send_message(
            chat_id=chat_id,
            text=escape(text),
            parse_mode="MarkdownV2",
        )

async def check_before_tarot(update: Update, context: ContextTypes.DEFAULT_TYPE):
    tarot_storage = TarotStorage()
    
    if await tarot_storage.is_today_drawn(update.message.chat_id):
        text="Already drawn today! Please check your chat history or use /history to check today's tarot history."
        await update.message.reply_text(
            text=escape(text),
            parse_mode="MarkdownV2",
        )
        return False
    chat_id = update.effective_chat.id
    result = TarotUser.get_user_info(str(chat_id))
    if result.get("status") == "success":
        if result.get("amount") < 20:
            text="You don't have enough points to draw a tarot. Please use /amount to check your current balance."
            await update.message.reply_text(
                text=escape(text),
                parse_mode="MarkdownV2",
            )
            return False
        else:
            return True
    return False
    
    