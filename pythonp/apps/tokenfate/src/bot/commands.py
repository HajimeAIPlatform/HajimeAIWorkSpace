from telegram import BotCommand
from telegram.ext import (
    CommandHandler,
    MessageHandler,
    ConversationHandler,
    ContextTypes,
    filters,
    Application
)
from os import getenv
from pythonp.apps.tokenfate.src.ton.views import send_tx, sell_transaction, buy_transaction, WAITING_FOR_INPUT, cancel
from pythonp.apps.tokenfate.src.bot.i18n_helper import I18nHelper
from pythonp.apps.tokenfate.models.transaction import UserPoints
from telegram import BotCommand, BotCommandScopeChat, BotCommandScopeDefault, Update

async def set_default_bot_commands(bot):
    """设置全局默认命令"""
    default_commands = [
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
    await bot.set_my_commands(commands=default_commands, scope=BotCommandScopeDefault())

async def set_user_specific_commands(bot, user_id, lang=None):
    """为特定用户设置个性化命令"""
    lang = lang or UserPoints.get_language_by_user_id(user_id) or 'zh'
    i18n = I18nHelper(lang)

    user_commands = [
        BotCommand("start", i18n.get_dialog("description_start")),
        BotCommand("quest", i18n.get_dialog("description_quest")),
        BotCommand("connect", i18n.get_dialog("description_connect")),
        BotCommand("disconnect", i18n.get_dialog("description_disconnect")),
        BotCommand("aura", i18n.get_dialog("description_aura")),
        BotCommand("language", i18n.get_dialog("description_language")),
    ]
    
    await bot.set_my_commands(
        commands=user_commands,
        scope=BotCommandScopeChat(chat_id=user_id)
    )

async def on_startup(application: Application):
    """Bot启动时设置默认命令"""
    await set_default_bot_commands(application.bot)

def setup_bot(application: Application):
    # 注册启动回调
    application.post_init = on_startup