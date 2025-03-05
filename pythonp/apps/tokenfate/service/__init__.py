
from flask import Blueprint
from pythonp.apps.tokenfate.service.bot.views import bot
from pythonp.apps.tokenfate.service.bot2.views import bot2
from pythonp.apps.tokenfate.service.bot_fortune_teller.views import bot_fortune_teller
from pythonp.apps.tokenfate.service.ton.web import ton
from pythonp.apps.tokenfate.service.dify.views import dify

from pythonp.apps.tokenfate.service.bot.views import run_bot1
from pythonp.apps.tokenfate.service.bot2.views import run_bot2
from pythonp.apps.tokenfate.service.bot_fortune_teller.views import run_bot_fortune_teller
import asyncio

blueprint = Blueprint('api', __name__, url_prefix='/telebot')
blueprint.register_blueprint(bot, url_prefix='/bot')
blueprint.register_blueprint(bot2, url_prefix='/bot2')
blueprint.register_blueprint(bot_fortune_teller, url_prefix='/bot_fortune_teller')
blueprint.register_blueprint(ton, url_prefix='/ton')
blueprint.register_blueprint(dify, url_prefix='/dify')

async def run_bot():
    await run_bot1()
    await run_bot2()
    await run_bot_fortune_teller()

loop = asyncio.get_event_loop()
loop.run_until_complete(run_bot())
