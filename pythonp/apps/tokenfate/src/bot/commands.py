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
from typing import List, Optional
from pythonp.apps.tokenfate.src.ton.views import send_tx, sell_transaction, buy_transaction, WAITING_FOR_INPUT, cancel
from pythonp.apps.tokenfate.src.bot.i18n_helper import I18nHelper
from pythonp.apps.tokenfate.models.transaction import UserPoints
from telegram import BotCommand, BotCommandScopeChat, BotCommandScopeDefault, Update

# 命令名称常量
CMD_START = "start"
CMD_QUEST = "quest"
CMD_CONNECT = "connect"
CMD_DISCONNECT = "disconnect"
CMD_AURA = "aura"
CMD_LANGUAGE = "language"

DEFAULT_COMMANDS = [
    (CMD_START, "description_start"),
    (CMD_QUEST, "description_quest"),
    (CMD_CONNECT, "description_connect"),
    (CMD_DISCONNECT, "description_disconnect"),
    (CMD_AURA, "description_aura"),
    (CMD_LANGUAGE, "description_language"),
]

def get_command_list(i18n, commands: List[tuple]) -> List[BotCommand]:
    """生成命令列表"""
    return [BotCommand(command, i18n.get_dialog(description)) for command, description in commands]

async def set_commands(bot, user_id: Optional[int] = None, lang: Optional[str] = None):
    """设置命令，支持全局命令和特定用户命令"""
    if user_id is not None:
        lang = lang or UserPoints.get_language_by_user_id(user_id) or 'zh'
        i18n = I18nHelper(lang)
        commands = get_command_list(i18n, DEFAULT_COMMANDS)
        try:
            await bot.set_my_commands(commands=commands, scope=BotCommandScopeChat(chat_id=user_id))
        except Exception as e:
            # 错误处理
            print(f"Failed to set commands for user {user_id}: {e}")
    else:
        # 全局命令默认使用中文
        commands = get_command_list(I18nHelper('zh'), DEFAULT_COMMANDS)
        await bot.set_my_commands(commands=commands)

# async def on_startup(application: Application):
#     """Bot启动时设置默认命令"""
#     await set_commands(application.bot)

# def setup_bot(application: Application):
#     # 注册启动回调
#     application.post_init = on_startup