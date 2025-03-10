'''
Description: 
Author: Devin
Date: 2025-03-06 11:04:12
'''

from flask import Blueprint
from pythonp.apps.fortune_teller.service.fortune_teller_bot.views import fortune_teller_bot
from pythonp.apps.fortune_teller.service.ton.web import ton
from pythonp.apps.fortune_teller.service.dify.views import dify

from pythonp.apps.fortune_teller.service.fortune_teller_bot.views import run_fortune_teller_bot
import asyncio

blueprint = Blueprint('api', __name__, url_prefix='/fortune_teller')
blueprint.register_blueprint(fortune_teller_bot, url_prefix='/bot_fortune_teller')
blueprint.register_blueprint(ton, url_prefix='/ton')
blueprint.register_blueprint(dify, url_prefix='/dify')

async def run_bot():
    await run_fortune_teller_bot()

loop = asyncio.get_event_loop()
loop.run_until_complete(run_bot())
