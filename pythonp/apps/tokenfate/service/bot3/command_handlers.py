from telegram import Update, ContextTypes

# 定义每个命令的处理逻辑
async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text("Welcome to the bot!")

async def history(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text("Here is your history.")

async def community(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text("Join our community!")

async def integral(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text("Your integral is 0.")