from flask import Blueprint
from pythonp.apps.tokenfate.service.bot.views import bot, run_bot
from pythonp.apps.tokenfate.service.bot2.views import bot2
from pythonp.apps.tokenfate.service.ton.web import ton
from pythonp.apps.tokenfate.service.dify.views import dify

# import src.bot.loop

blueprint = Blueprint('api', __name__, url_prefix='/telebot')
blueprint.register_blueprint(bot, url_prefix='/bot')
blueprint.register_blueprint(bot2, url_prefix='/bot2')
blueprint.register_blueprint(ton, url_prefix='/ton')
blueprint.register_blueprint(dify, url_prefix='/dify')
