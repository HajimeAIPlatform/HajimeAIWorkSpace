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
    ("tarot", "description_tarot"),
    ("history", "description_quest"),
    ("community", "description_connect"),
    ("amount", "despcription_integral"),
]

def get_command_list(commands: List[tuple]) -> List[BotCommand]:
    """生成命令列表"""
    print("Generating command list...")
    return [BotCommand(command, "loading...") for command, description in commands]

async def set_commands(bot):
    """设置命令，支持全局命令和特定用户命令"""
    print("Setting commands...")
    commands = get_command_list(DEFAULT_COMMANDS)
    await bot.set_my_commands(commands=commands)
