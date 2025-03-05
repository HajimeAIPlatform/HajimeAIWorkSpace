from telegram import BotCommand


async def set_bot_commands(telegram_app):
    commands = [
        BotCommand("start", "Start Use Fortune Teller"),
    ]
    await telegram_app.bot.set_my_commands(commands)

async def set_bot_commands_handler(telegram_app):
    await set_bot_commands(telegram_app)
