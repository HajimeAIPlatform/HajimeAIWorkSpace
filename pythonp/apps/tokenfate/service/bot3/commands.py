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
from pythonp.apps.tokenfate.service.ton.views import send_tx, sell_transaction, buy_transaction, WAITING_FOR_INPUT, cancel
from pythonp.apps.tokenfate.models.transaction import UserPoints
from telegram import BotCommand, BotCommandScopeChat, BotCommandScopeDefault, Update

DEFAULT_COMMANDS = [
    # example: ("command", "description")
    ("start", "description_start"),
    ("history", "description_quest"),
    ("community", "description_connect"),
    ("integral", "despcription_integral"),
]

def get_command_list(i18n, commands: List[tuple]) -> List[BotCommand]:
    """生成命令列表"""
    return [BotCommand(command, "loading...") for command, description in commands]

async def set_commands(bot, user_id: Optional[int] = None, lang: Optional[str] = None):
    """设置命令，支持全局命令和特定用户命令"""
    print("Setting commands...")
    lang = lang or UserPoints.get_language_by_user_id(user_id) or 'zh'
    i18n = None
    if user_id:
        commands = get_command_list(i18n, DEFAULT_COMMANDS)
        try:
            await bot.set_my_commands(commands=commands, scope=BotCommandScopeChat(chat_id=user_id))
        except Exception as e:
            # 错误处理
            print(f"Failed to set commands for user {user_id}: {e}")
    else:
        # 全局命令默认使用中文
        commands = get_command_list(i18n, DEFAULT_COMMANDS)
        await bot.set_my_commands(commands=commands)
