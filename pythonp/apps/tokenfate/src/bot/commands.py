from telegram import BotCommand,InlineKeyboardButton, InlineKeyboardMarkup, WebAppInfo,MenuButtonWebApp
from telegram.ext import (
    CommandHandler,
    MessageHandler,
    ConversationHandler,
    ContextTypes,
    filters,
)
from os import getenv
from pythonp.apps.tokenfate.src.ton.views import send_tx, sell_transaction, buy_transaction, WAITING_FOR_INPUT, cancel
from pythonp.apps.tokenfate.src.bot.i18n_helper import I18nHelper

WEB_MINI_APP_URL = getenv('WEB_MINI_APP_URL')

i18n = I18nHelper()
async def set_menu_button(telegram_app):
    # 创建一个MenuButtonWebApp对象
    miniapp_url = f"{WEB_MINI_APP_URL}"
    web_app_info = WebAppInfo(url=miniapp_url)
    menu_button = MenuButtonWebApp(text="Wallet", web_app=web_app_info)

    # 设置主菜单按钮
    await telegram_app.bot.set_chat_menu_button(menu_button=menu_button)


async def set_bot_commands(telegram_app):
    commands = [
        BotCommand("start", "开始使用 TokenFate"),
        BotCommand("quest", "探问符命"),
        BotCommand("connect", "连接钱包"),
        # BotCommand("buy", "Buy transaction"),
        # BotCommand("sell", "Sell transaction"),
        BotCommand("disconnect", "断开钱包连接"),
        # BotCommand("my_wallet", "Show connected wallet"),
        BotCommand("aura", "符命，灵能感应！"),
        BotCommand("language", "语言设置"),
    ]
    await telegram_app.bot.set_my_commands(commands)
    # await set_menu_button(telegram_app)


async def set_bot_send_ex_handler(telegram_app):
    buy_handler = ConversationHandler(
        entry_points=[CommandHandler('buy', send_tx)],
        states={
            WAITING_FOR_INPUT:
            [MessageHandler(filters.TEXT & ~filters.COMMAND, buy_transaction)],
        },
        fallbacks=[CommandHandler('cancel', cancel)],
    )

    sell_handler = ConversationHandler(
        entry_points=[CommandHandler('sell', send_tx)],
        states={
            WAITING_FOR_INPUT: [
                MessageHandler(filters.TEXT & ~filters.COMMAND,
                               sell_transaction)
            ],
        },
        fallbacks=[CommandHandler('cancel', cancel)],
    )

    telegram_app.add_handler(buy_handler)
    telegram_app.add_handler(sell_handler)


async def set_bot_commands_handler(telegram_app):
    await set_bot_commands(telegram_app)
