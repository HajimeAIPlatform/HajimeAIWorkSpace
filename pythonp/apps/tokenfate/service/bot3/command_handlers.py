from md2tgmd import escape
from pythonp.apps.tokenfate.service.dify.views import chat_tarot
from telegram import (Update, helpers)

from telegram.ext import (
    ContextTypes,
)

# 定义每个命令的处理逻辑
async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text("Welcome to the bot!")

# async def tarot(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
#     # 获取参数
#     args = context.args
#     text="示例: /tarot question"
#     if not args:
#         await context.bot.send_message(chat_id=update.effective_chat.id, text=escape(text), parse_mode="MarkdownV2")
#         return
    
#     # 根据不同参数处理不同的塔罗牌逻辑
#     inputs = {
#         "input": args,
#     }
#     answer = chat_tarot(inputs)
#     await context.bot.send_photo(
#         photo=answer["url"],
#         caption=escape(answer["text"]),
#         parse_mode="MarkdownV2",
#     )

async def tarot(update: Update, context: ContextTypes.DEFAULT_TYPE):
    # 获取用户输入的参数
    args = context.args
    chat_id = update.effective_chat.id,
    if not args:
        # 如果没有参数，提示用户 /tarot 的使用方法
        text="Usage: /tarot question"
        await update.message.reply_text(
            text=escape(text),
            parse_mode="MarkdownV2",
        )
    else:
        # 如果有参数，执行 chat_tarot 处理逻辑
        user_question = " ".join(args)  # 合并所有参数为完整问题
        inputs = {
            "input": user_question,
        }
        answer = chat_tarot(inputs)
        await update.message.reply_photo(
            photo=answer["url"],
            caption=escape(answer["text"]),
            parse_mode="MarkdownV2",
        )

async def history(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text("Here is your history.")

async def community(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text("Join our community!")

async def integral(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    await update.message.reply_text("Your integral is 0.")
