import logging
import asyncio
from io import BytesIO
from os import getenv
import urllib.parse
import time
import json
import os
import re
from typing import List, Dict, Union

from flask import Blueprint, jsonify, request
from md2tgmd import escape
from telegram import Update, BotCommand, InlineKeyboardButton, InlineKeyboardMarkup, WebAppInfo, \
    InlineQueryResultsButton, InputMediaPhoto, Message, CallbackQuery
from telegram.ext import ApplicationBuilder, DictPersistence, CommandHandler

from pythonp.apps.tokenfate.src.dify.views import chat_blocking, chat_streaming, chat_workflow, chat_decode
from pythonp.apps.tokenfate.src.binance.views import handle_binance_command
from pythonp.apps.tokenfate.src.binance.utils import get_all_prices, process_recommendation
import pythonp.apps.tokenfate.src.ton.views as ton_module
from pythonp.apps.tokenfate.src.bot.commands import set_commands
from pythonp.apps.tokenfate.src.bot.wallet_menu_callback import set_handlers
import pythonp.apps.tokenfate.src.bot.state as ChatStatus
from pythonp.apps.tokenfate.src.bot.i18n_helper import I18nHelper
from pythonp.apps.tokenfate.src.bot.keyboards import KeyboardFactory
from pythonp.apps.tokenfate.src.ton.tc_storage import DailyFortune, UserActivityTracker
from pythonp.apps.tokenfate.models.transaction import UserPoints
from pythonp.apps.tokenfate.static.static import get_images_path
from pythonp.common.email_notifications.error_collector import collect_error

# 获取Telegram Bot Token
telegram_bot_token = getenv('TELEGRAM_BOT_TOKEN')
if not telegram_bot_token:
    logging.error("Telegram bot token is not set in the environment")
    raise ValueError("Telegram bot token is not set in the environment")

# 创建Telegram应用
persistence = DictPersistence()
telegram_app = ApplicationBuilder().token(telegram_bot_token).persistence(
    persistence).build()

# 初始化redis数据库连接
daily_fortune = DailyFortune()
user_activity_tracker = UserActivityTracker()

WEB_MINI_APP_URL = getenv('WEB_MINI_APP_URL')


async def run_bot():
    await set_commands(telegram_app.bot)
    # await set_bot_commands_handler(telegram_app)
    await telegram_app.initialize()
    await telegram_app.start()

loop = asyncio.get_event_loop()
loop.run_until_complete(run_bot())

# 创建Flask Blueprint
bot = Blueprint('bot', __name__)


async def handle_streaming_chat(chat_id, query):
    data = {"query": query}
    generator = chat_streaming(data)

    try:
        # 先发送一条初始消息，并获取 message_id
        initial_message = await telegram_app.bot.send_message(
            chat_id=chat_id,
            text="Generating content, please wait...",
        )
        message_id = initial_message.message_id
        current_message_content = initial_message.text
        message = ""
        # 使用 message_id 更新消息内容
        for chunk in generator:
            if chunk == "" or chunk is None: continue
            message = message + chunk
            escaped_chunk = escape(message)
            logging.debug(f"Escaped chunk: {escaped_chunk}")

            # 仅在新内容与当前内容不同时更新消息
            if message != current_message_content:
                await telegram_app.bot.edit_message_text(
                    chat_id=chat_id,
                    message_id=message_id,
                    text=escaped_chunk,
                    parse_mode="MarkdownV2"
                )
                current_message_content = message

    except Exception as e:
        logging.error(f"Error occurred during streaming: {e}")
        collect_error(e)
        await telegram_app.bot.send_message(
            chat_id=chat_id,
            text='Sorry, I am not able to generate content for you right now. Please try again later.'
        )


async def ton_command_handle(update):
    try:
        ton_response = await ton_module.handle_ton_command(
            telegram_app, update)
        if ton_response:
            return ton_response
        return None

    except Exception as e:
        logging.error(f"Error occurred during streaming: {e}")
        collect_error(e)
        return {
            "text":
                'Sorry, I am not able to generate content for you right now. Please try again later.'
        }


@bot.route('/binance/price', methods=['GET'])
async def binance_price():
    return get_all_prices(10000, "USDT")


# async def start(update):
#     # 创建没有参数的按钮
#     web_app_info = WebAppInfo(url=WEB_MINI_APP_URL + "/#/", api_kwargs={"aa": "11"})
#     button = InlineQueryResultsButton(text="Open Wallet", web_app=web_app_info)

#     # 创建带有参数的按钮
#     action = "buy"
#     token = "ORN"
#     amount = 100
#     timestamp = int(time.time())  # 添加时间戳
#     web_app_url_with_params = f"{WEB_MINI_APP_URL}/#/transaction?action={urllib.parse.quote(action)}&token={urllib.parse.quote(token)}&amount={urllib.parse.quote(str(amount))}&t={timestamp}"
#     web_app_info2 = WebAppInfo(url=web_app_url_with_params, api_kwargs={"aa": "22"})
#     button2 = InlineQueryResultsButton(text="111 Wallet", web_app=web_app_info2)

#     keyboard = InlineKeyboardMarkup([[button], [button2]])

#     await update.message.reply_text("Click a button to open Web App", reply_markup=keyboard)


@bot.route('/webhook', methods=['POST'])
# @run_async
async def webhook():
    chat_id = None
    try:
        body = request.get_json()
        update = Update.de_json(body, telegram_app.bot)
        logging.debug(f"Received update: {update}")
        
        if update.edited_message:
            return jsonify({'status': 'ok'}), 200
        # 处理update.message和update.callback_query.message
        chat_id = await get_chat_id(update) # 获取chat_id
        print(UserPoints)
        lang = UserPoints.get_language_by_user_id(chat_id)

        if update.callback_query and update.callback_query.data.startswith("lang"):
            data = update.callback_query.data
            logging.info(f"data: {data}")
            details = data.split(":")
            lang = details[1]
            await set_commands(telegram_app.bot, chat_id, lang)
            await update_default_language(update, lang=lang)
            return jsonify({'status': 'ok'}), 200

        if not lang:
            await set_language(update)
            return jsonify({'status': 'ok'}), 200
        
        # 初始化语言和键盘工厂
        i18n = I18nHelper(lang)
        keyboard_factory = KeyboardFactory(i18n)

        await handle_daily_checkin(update)
        if update.callback_query:
            await set_handlers(update, telegram_app)
            user_id = update.callback_query['from']['id']
            if update.callback_query.data == "for_your_information_button":
                await update.callback_query.answer()
                
                # 获取对话文本
                dialog = i18n.get_dialog('info')
                
                # 创建新的键盘布局
                reply_markup = keyboard_factory.create_keyboard("info")
                
                # 发送新的消息
                await update.callback_query.message.reply_text(
                    escape(dialog),
                    parse_mode="MarkdownV2",
                    reply_markup=reply_markup
                )
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data == "launch_to_reveal_button":
                await update.callback_query.answer()
                
                # 获取对话文本
                dialog = i18n.get_dialog('quest')
                
                # 发送新的消息
                await update.callback_query.message.reply_text(
                    escape(dialog),
                    parse_mode='MarkdownV2'
                )
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data == "show_aura_rules":
                await update.callback_query.answer()
                await show_aura_rules(update)
                return jsonify({'status': 'ok'}), 200

            if update.callback_query.data == "connect_wallet_button":
                await update.callback_query.answer()
                # connect to wallet
                connected_wallets = await user_activity_tracker.get_connected_wallets(user_id)
                
                if connected_wallets:
                    # Format the set into a list with each wallet on a new line
                    wallet_list = "\n".join(connected_wallets)
                    dialog = i18n.get_dialog("connected_wallets")
                    dialog = dialog.format(wallet_list=wallet_list)
                    await update.callback_query.message.reply_text(escape(dialog), parse_mode="MarkdownV2")
                await ton_module.wallet_menu_callback.on_choose_wallet_click(update)
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data == "aura_action_daily_checkin":
                await update.callback_query.answer()
                if await user_activity_tracker.is_checked_in_today(user_id=user_id):
                    await update.callback_query.message.reply_text(i18n.get_dialog("aura_action_daily_checkin_again"))
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data == "aura_action_recommend_click":
                await update.callback_query.answer()
                # await send_recommendations(update)
                await risk_preference(update, token=None, risk_target="recommend_only")
                return jsonify({'status': 'ok'}), 200

            if update.callback_query.data.startswith("reveal_fate"):
                data = update.callback_query.data
                logging.info(f"data: {data}")
                details = data.split(":")
                token = details[1]
                await risk_preference(update, token, risk_target="reveal_fate")
                return jsonify({'status': 'ok'}), 200

            if update.callback_query.data.startswith("recommend"):
                data = update.callback_query.data
                logging.info(f"data: {data}")
                details = data.split(":")
                token = details[1]
                if handle_recommendation_click(chat_id):
                    await get_aura_status(update, "aura_action_recommend_click")
                await risk_preference(update, token, risk_target="reveal_fate")
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data.startswith("risk"):
                callback_data = update.callback_query.data
                logging.info(f"callback_data: {callback_data}")

                # 先去掉前缀 "risk:" 并获取剩余部分
                remaining_data = callback_data.split(':')[1]
                
                # 使用 & 分割字符串，得到一个列表，每个元素都是 key=value 形式的字符串
                details = remaining_data.split("&")
                
                # 创建一个字典来存储解析出来的键值对
                params = {}
                for detail in details:
                    # 每个 detail 是 key=value 形式，再次使用 = 分割
                    key, value = detail.split("=")
                    params[key] = value
                
                # 从 params 字典中获取 token, risk_target 和 role 的值
                token = params.get('token')
                risk_target = params.get('risk_target')
                role = params.get('role')
                
                # 打印或记录获取到的值，以便调试
                logging.info(f"token: {token}, risk_target: {risk_target}, role: {role}")
                
                # 调用函数，假设函数需要 token 和 role 参数
                if risk_target == "reveal_fate":
                    await reveal_fate(update, token, role)  # 如果 reveal_fate 需要 role 参数，则传入
                elif risk_target == "recommend_only":
                    await send_recommendations(update, role)
                
                # 返回 HTTP 响应
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data.startswith("decode"):
                data = update.callback_query.data
                logging.info(f"data: {data}")
                details = data.split(":")
                status, token, role = details
                if status == "decode_yes":
                    await decode_lot(update, token, role)
                elif status == "decode_no":
                    await message_unveil_or_not(update, token, role)
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data.startswith("menu"):
                await message_menu(update)
                return jsonify({'status': 'ok'}), 200

        if update.message:
        # Process update with the application
            await telegram_app.process_update(update)
            result = ChatStatus.get_transaction_status(update.message.chat_id)
            print(result, 'result', update.message.chat_id)
            if result:
                return jsonify({'status': 'ok'}), 200
            if update.message.text == '/cancel' or update.message.text == '/buy' or update.message.text == '/sell':
                return jsonify({'status': 'ok'}), 200

            ton_response = await ton_command_handle(update)
            if ton_response:
                if isinstance(ton_response, dict) and "text" in ton_response:
                    return {
                        "method": "sendMessage",
                        "chat_id": update.message.chat_id,
                        "text": ton_response["text"],
                    }
                print(ton_response, 'return')
                return jsonify({'status': 'ok'}), 200

            chat_id = update.message.chat_id
            binance_response = handle_binance_command(update.message.text)
            if binance_response:
                return {
                    "method": "sendMessage",
                    "chat_id": chat_id,
                    "text": binance_response,
                }

            if update.message.text == '/start':
                await start(update)
                return jsonify({'status': 'ok'}), 200

            if update.message.text == '/quest':
                dialog = i18n.get_dialog('quest')
                await update.message.reply_text(
                    text=escape(dialog),
                    parse_mode='MarkdownV2'
                )
                return jsonify({'status': 'ok'}), 200
            
            if update.message.text == '/aura':
                await show_aura_rules(update)
                return jsonify({'status': 'ok'}), 200

            if update.message.text == '/language':
                await set_language(update)
                return jsonify({'status': 'ok'}), 200

            if update.message.text.startswith('$') and validate_ticker(update.message.text[1:]):
                chat_id = update.message.chat_id
                Token = update.message.text
                print(Token, 'Token')
                # 当收到以 $ 开头的消息时，发送新的消息并附带 Connect Wallet 按钮
                is_connected = await ton_module.check_connected(update, telegram_app)
                
                if not is_connected:
                    # 临时保存 Token
                    ChatStatus.set_tmp_token(chat_id, Token)
                    
                    # 获取未连接钱包时的对话文本
                    dialog = i18n.get_dialog('unconnected')                
                    # 创建按钮：连接钱包或直接揭示
                    reply_markup = keyboard_factory.create_keyboard("unconnected", token=Token)
                    await update.message.reply_text(escape(dialog), parse_mode="MarkdownV2", reply_markup=reply_markup)
                else:
                    # 获取已连接钱包时的对话文本
                    dialog = i18n.get_dialog('connected')                
                    # 创建新的按钮：揭示代币命运
                    reply_markup = keyboard_factory.create_keyboard("connected", token=Token)
                    
                    # 发送新的消息
                    await update.message.reply_text(escape(dialog), parse_mode="MarkdownV2", reply_markup=reply_markup)
                return jsonify({'status': 'ok'}), 200

            if update.message.photo:
                file_id = update.message.photo[-1].file_id
                logging.info(f"Images file id is {file_id}")
                file = await telegram_app.bot.get_file(file_id)
                logging.info("Image file found")
                bytes_array = await file.download_as_bytearray()
                bytesIO = BytesIO(bytes_array)
                logging.info("Images file as bytes")
                image = Image.open(bytesIO)
                logging.info("Image opened")

                prompt = 'Describe the image'

                if update.message.caption:
                    prompt = update.message.caption
                logging.info(f"Prompt is {prompt}")

                text = "test"

                return {
                    "method": "sendMessage",
                    "chat_id": chat_id,
                    "text": escape(text),
                    "parse_mode": "MarkdownV2"
                }
            else:
                api_responseponse = chat_blocking({
                    "query": update.message.text,
                    "user": chat_id
                })
                # await handle_streaming_chat(chat_id, update.message.text)
                text = api_responseponse
                result = process_recommendation(api_responseponse)
                if isinstance(result, tuple):
                    token, action, amount = result
                    # 创建按钮
                    button_text = f"{action} {token}:{amount}"

                    # 构建小程序的 URL
                    miniapp_url = f"{WEB_MINI_APP_URL}/transaction?action={urllib.parse.quote(action)}&token={urllib.parse.quote(token)}&amount={urllib.parse.quote(str(amount))}"
                    print(miniapp_url, 'miniapp_url')
                    web_app_info = WebAppInfo(url=miniapp_url)
                    button = InlineKeyboardButton(text=f"{action.upper()} {token}:{amount}", web_app=web_app_info)
                    keyboard = InlineKeyboardMarkup([[button]])

                    await update.message.reply_text(escape(text), parse_mode="MarkdownV2", reply_markup=keyboard)
                else:
                    await update.message.reply_text(escape(text), parse_mode="MarkdownV2")
                return jsonify({'status': 'ok'}), 200
        return jsonify({'status': 'ok'}), 200

    except Exception as e:
        logging.info(f"update: {update}")
        logging.error(f"Error Occurred: {e}")
        collect_error(e)
        return {
            "method":
                "sendMessage",
            "chat_id":
                chat_id,
            "text":
                'Sorry, I am not able to generate content for you right now. Please try again later.'
        }
        # return 'OK'

def validate_token_data(response_data: Dict[str, Union[str, list]]) -> List[Dict[str, str]]:
    if not isinstance(response_data, dict) or 'data' not in response_data:
        logging.error("Invalid response format: missing 'data' field")
        return []

    tokens = response_data['data']
    if not isinstance(tokens, list):
        logging.error("Invalid tokens format: 'data' is not a list")
        return []

    # Additional validation of token data structure
    validated_tokens = []
    for token_info in tokens:
        if isinstance(token_info, dict) and 'token' in token_info and 'reason' in token_info:
            validated_tokens.append({
                'token': str(token_info['token']).upper(),
                'reason': str(token_info['reason'])
            })

    return validated_tokens


def parse_token_response(response: str) -> Dict[str, Union[str, list]]:
    try:
        return json.loads(response)
    except json.JSONDecodeError as e:
        logging.error("Failed to parse API response: %s", e)
        raise


def fetch_trending_tokens(role: str):
    try:
        # Fetch recommendations from Dify API
        data = {
            "risk_preference": role,
        }
        api_response = chat_workflow(data)
        logging.info("Raw API response: %s", api_response)
        parsed_response = parse_token_response(api_response)
        logging.info(f"parsed_response: {parsed_response}")
        return validate_token_data(parsed_response)

    except Exception as e:
        logging.error("Failed to fetch trending tokens: %s", e)
        collect_error(e)
        return []
    
def create_token_keyboard(tokens):
    keyboard = []
    row = []

    for index, token_info in enumerate(tokens):
        token_name = token_info.get('token', '')
        if token_name and index < 4:  # 只处理前四个token
            button = InlineKeyboardButton(
                text=f"${token_name}",
                callback_data=f"recommend:${token_name}"
            )
            row.append(button)

            if len(row) == 2:  # 每行2个按钮
                keyboard.append(row)
                row = []

    return InlineKeyboardMarkup(keyboard)


async def send_recommendations(update, role: str):
    if update.message:
        target = update
        user_id = update.message['from']['id']
    # 处理回调查询
    elif update.callback_query:
        target = update.callback_query
        user_id = update.callback_query['from']['id']
    else:
        return
    lang = UserPoints.get_language_by_user_id(user_id)
    i18n = I18nHelper(lang)

    """发送推荐的tokens给用户"""
    recommended_tokens = fetch_trending_tokens(role)
    if not recommended_tokens:
        logging.error("Failed to get token recommendations")
        return

    reply_markup = create_token_keyboard(recommended_tokens)
    dialog = i18n.get_dialog("unveil_result")
    await target.message.reply_text(escape(dialog), parse_mode="MarkdownV2", reply_markup=reply_markup)


async def reveal_fate(update, token, role: str):
    try:
        # 处理普通消息
        if update.message:
            target = update
            user_id = update.message['from']['id']
        # 处理回调查询
        elif update.callback_query:
            target = update.callback_query
            user_id = update.callback_query['from']['id']
        else:
            return
        lang = UserPoints.get_language_by_user_id(user_id)
        i18n = I18nHelper(lang)
        keyboard_factory = KeyboardFactory(i18n)

        if not token:
            await target.message.reply_text("请输入你的符币代码。")
            return jsonify({'status': 'ok'}), 200

        # 随机抽签
        cached_lot = await daily_fortune.get_cached_lot(user_id, token)
        is_cached = cached_lot is not None
        if is_cached:
            result_of_draw = cached_lot
        else:
            result_of_draw = await daily_fortune.get_daily_lot(user_id, token)
        logging.info(f"result_of_draw: {result_of_draw}")

        if is_cached:
            # 用户积分无需更新
            await get_aura_status(update, "aura_action_cached_hit")
        else:
            sql_status = UserPoints.update_points_by_user_id(user_id=user_id, points=-5)
            if not sql_status:
                logging.error("Failed to update user points")
                await get_aura_status(update, "aura_action_invalid")
            await get_aura_status(update, "aura_action_fate_reveal")

        sign_level = result_of_draw["sign_level"]
        sign_from = result_of_draw["sign_from"]
        sign_text = result_of_draw["sign_text"]
        # 发送抽签结果
        image_path = get_images_path(f'{sign_level}.png')
        dialog = i18n.get_dialog("lot_daily_content")
        dialog = dialog.format(token=token, sign_from=sign_from, sign_text=sign_text)
        reply_markup = keyboard_factory.create_keyboard("lot", token=token, role=role)
        with open(image_path, 'rb') as image_file:
            await target.message.reply_photo(
                photo=image_file,
                caption=escape(dialog),
                parse_mode="MarkdownV2",
                reply_markup=reply_markup,
            )
        # 发送推荐的tokens
        # await send_recommendations(update, role)

        return jsonify({'status': 'ok'}), 200
    
    except Exception as e:
        logging.error(f"Error in reveal_fate: {e}")
        collect_error(e)
        await target.message.reply_text(
            "Sorry, something went wrong while revealing your fate."
        )
        return

    
async def risk_preference(update, token, risk_target):
    chat_id = await get_chat_id(update) # 获取chat_id
    lang = UserPoints.get_language_by_user_id(chat_id)
    i18n = I18nHelper(lang)
    try:
        keyboard_factory = KeyboardFactory(i18n)
        await update.callback_query.answer()
        # if not token:
        #     await update.callback_query.message.reply_text("Please Enter Your Token.")
        #     return jsonify({'status': 'ok'}), 200
        
        # 发送文本信息
        dialog = i18n.get_dialog('risk')
        await update.callback_query.message.reply_text(
            text=dialog,
            parse_mode="MarkdownV2"
        )

        # 根据触发源选择键盘
        if risk_target == "reveal_fate":
            reply_markup = keyboard_factory.create_keyboard("risk", token=token, risk_target="reveal_fate")
        elif risk_target == "recommend_only":
            reply_markup = keyboard_factory.create_keyboard("risk", token=token, risk_target="recommend_only")
        else:
            reply_markup = keyboard_factory.create_keyboard("risk", token=token)

        # 发送图片及选择项
        image_path = get_images_path('risk_preference_combined.png')
        with open(image_path, 'rb') as image_file:
            await update.callback_query.message.reply_photo(
                photo=image_file,
                parse_mode="MarkdownV2",
                reply_markup=reply_markup
            )
        return jsonify({'status': 'ok'}), 200

    except Exception as e:
        logging.error(f"Error in reveal_fate: {e}")
        collect_error(e)
        await update.callback_query.message.reply_text(
            "Sorry, something went wrong while showing risk preference."
        )
        return

    
async def show_aura_rules(update):    
    try:
        if update.callback_query:
            await update.callback_query.answer()
            target = update.callback_query.message
            user_id = update.callback_query['from']['id']
        elif update.message:
            target = update.message
            user_id = update.message['from']['id']
        else:
            return
        lang = UserPoints.get_language_by_user_id(user_id)
        i18n = I18nHelper(lang)
        keyboard_factory = KeyboardFactory(i18n)
        # Send aura rules information
        points = UserPoints.get_points_by_user_id(user_id=user_id)
        daily_recommended_points = UserPoints.get_daily_recommended_points(user_id=user_id)
        aura_status_amount = i18n.get_dialog('aura_status_amount').format(points=points)
        aura_status_daily = i18n.get_dialog('aura_status_daily').format(daily_recommended_points=daily_recommended_points)
        aura_rules = i18n.get_dialog('aura_rules')
        dialog = aura_status_amount + aura_status_daily + aura_rules
        # logging.info(dialog)
        await target.reply_text(
            text = escape(dialog),
            parse_mode="MarkdownV2"
        )

        # Send interactive elements
        reply_markup = keyboard_factory.create_keyboard("aura")
        image_name = f'{lang}-ways-to-impact-aura.png'
        image_path = get_images_path(image_name)
        with open(image_path, 'rb') as image_file:
            await target.reply_photo(
                photo=image_file,
                reply_markup=reply_markup
            )
        return jsonify({'status': 'ok'}), 200

    except Exception as e:
        logging.error(f"Error in show_aura_rules: {e}")
        collect_error(e)
        await target.reply_text("Sorry, something went wrong while showing the aura rules.")
        return

    
async def decode_lot(update, token, role):
    try:
        # 处理普通消息
        if update.message:
            target = update
            user_id = update.message['from']['id']
        # 处理回调查询
        elif update.callback_query:
            target = update.callback_query
            user_id = update.callback_query['from']['id']
        else:
            return
        
        lang = UserPoints.get_language_by_user_id(user_id)
        i18n = I18nHelper(lang)
        keyboard_factory = KeyboardFactory(i18n)
        # 获取今日缓存的签
        result_of_draw = await daily_fortune.get_cached_lot(user_id, token)
        if result_of_draw is None:
            await target.message.reply_text("昨日之签无法解读。")
            return
        # 检查是否获取过签解不论语言，进而决定是否扣除积分
        decode_status = await daily_fortune.get_decode_status(user_id, token)
        if not decode_status:
            await get_aura_status(update, "aura_action_decode")

        # 先有签再有签解，获取今日缓存的签解
        cached_decode = await daily_fortune.get_cached_decode(user_id, token, lang)
        if cached_decode is None:
            # 签解为空，进行签解并缓存
            sign_text = result_of_draw['sign_text']
            data = {
                "query": sign_text,
                "inputs": {
                    "lot": sign_text,
                    "lang": lang,
                    "risk_preference": role,
                },
            }
            cached_decode = chat_decode(data)
            await daily_fortune.set_decode_cache(user_id, token, lang, cached_decode)

        # 发送签解
        dialog = i18n.get_dialog("lot_decoded_content")
        dialog = dialog.format(token=token, cached_decode=cached_decode)
        reply_markup = keyboard_factory.create_keyboard("decode")
        await target.message.reply_text(
            escape(dialog), parse_mode="MarkdownV2", 
            reply_markup=reply_markup
        )
        return jsonify({'status': 'ok'}), 200

    except Exception as e:
        logging.error(f"Error in decode_lot: {e}")
        collect_error(e)
        await target.message.reply_text("Sorry, something went wrong while decode your lot.")
        return

    
async def get_aura_status(update, action: str):
    try:
        # 处理普通消息
        if update.message:
            target = update
            user_id = update.message['from']['id']
        # 处理回调查询
        elif update.callback_query:
            target = update.callback_query
            user_id = update.callback_query['from']['id']
        else:
            return
        lang = UserPoints.get_language_by_user_id(user_id)
        i18n = I18nHelper(lang)
        dialog = i18n.get_dialog(action)
        dialog = dialog.format(user_id=user_id)
        button = i18n.get_button("aura")
        reply_markup = InlineKeyboardMarkup([[InlineKeyboardButton(text=button, callback_data="show_aura_rules")]])
        await target.message.reply_text(escape(dialog), parse_mode="MarkdownV2", reply_markup=reply_markup)
        return jsonify({'status': 'ok'}), 200

    except Exception as e:
        logging.error(f"Error in get_aura_status: {e}")
        collect_error(e)
        await target.message.reply_text("Sorry, something went wrong while getting your aura status.")
        return
    

async def handle_daily_checkin(update):
    try:
        # 获取用户ID
        if update.message:
            user_id = update.message['from']['id']
        elif update.callback_query:
            user_id = update.callback_query['from']['id']
        else:
            return

        # 检查用户是否已打卡
        if await user_activity_tracker.is_checked_in_today(user_id=user_id):
            return

        # 执行打卡操作
        if await user_activity_tracker.daily_checkin(user_id=user_id):
            # 更新用户积分
            sql_status = UserPoints.update_points_by_user_id(user_id=user_id, points=20)
            if not sql_status:
                logging.error("Failed to update user points")
                await get_aura_status(update, "aura_action_invalid")
                return

            # 发送打卡成功消息
            await get_aura_status(update, "aura_action_daily_checkin")
            logging.info(f"User {user_id} checked in successfully")
        else:
            logging.error("User {user_id} failed to check in")
            

    except Exception as e:
        logging.error(f"Error in handle_daily_checkin: {str(e)}")
        collect_error(e)
        await get_aura_status(update, "aura_action_invalid")


async def update_default_language(update, lang: str):
    try:
        if update.message:
            target = update
            user_id = update.message['from']['id']
        elif update.callback_query:
            target = update.callback_query
            user_id = update.callback_query['from']['id']
        else:
            return
        UserPoints.update_language_by_user_id(user_id=user_id, language=lang)
        await start(update)
        return 
    except Exception as e:
        logging.error(f"Error in update_default_language: {e}")
        collect_error(e)
        await target.message.reply_text("Sorry, something went wrong while setting your language.")
        return

    
async def set_language(update, i18n=I18nHelper()):
    keyboard_factory = KeyboardFactory(i18n)
    dialog = i18n.get_dialog('setting_lang')
    reply_markup = keyboard_factory.create_keyboard("lang")
    await update.message.reply_text(
        text=escape(dialog),
        parse_mode="MarkdownV2",
        reply_markup=reply_markup
    )

async def get_chat_id(update):
    try:
        if update.message:
            user_id = update.message['from']['id']
        elif update.callback_query:
            user_id = update.callback_query['from']['id']
        else:
            return
        return user_id
    except Exception as e:
        logging.error(f"Error in get_chat_id: {e}")
        collect_error(e)
        # await target.message.reply_text("Sorry, something went wrong while getting your chat ID.")
        return
    

async def start(update):
    try:
        # 处理普通消息
        if update.message:
            target = update
            user_id = update.message['from']['id']
        # 处理回调查询
        elif update.callback_query:
            target = update.callback_query
            user_id = update.callback_query['from']['id']
        else:
            return
        lang = UserPoints.get_language_by_user_id(user_id)
        i18n = I18nHelper(lang)
        keyboard_factory = KeyboardFactory(i18n)
        # 创建键盘
        reply_markup = keyboard_factory.create_keyboard("start")
        
        # 获取对话内容
        dialog = i18n.get_dialog('start')
        
        # 获取图片路径
        image_path = get_images_path(f"{lang}_welcome.png")

        # 检查文件是否存在
        if not os.path.exists(image_path):
            logging.error(f"Image file {image_path} does not exist.")
            await target.message.reply_text(text="There was an error loading the welcome image.", reply_markup=reply_markup)
            return ('Error', 500)
        
        # 发送图片和对话
        with open(image_path, 'rb') as image_file:
            await target.message.reply_photo(
                photo=image_file,
                caption=escape(dialog),
                parse_mode="MarkdownV2",
                reply_markup=reply_markup
            )
        return ('OK', 200)
    
    except Exception as e:
        logging.error(f"An error occurred in the start command: {e}")
        collect_error(e)
        await target.message.reply_text(text="An unexpected error occurred.")
        return ('Error', 500)
    
def handle_recommendation_click(user_id):
    points_to_add = 10
    if UserPoints.update_daily_recommended_points(user_id, points_to_add):
        if UserPoints.update_points_by_user_id(user_id, points_to_add, description="Recommendation click"):
            logging.info(f"User {user_id} successfully added {points_to_add} points for recommendation click, and successed toupdate user points.")
            return True
        else:
            logging.info(f"User {user_id} successfully added {points_to_add} points for recommendation click, but failed to update user points.")
            return False
    else:
        logging.info(f"User {user_id} has already reached the daily limit for recommendation clicks.")
        return False

async def message_menu(update):
    try:
        if update.message:
            target = update
            user_id = update.message['from']['id']
        elif update.callback_query:
            target = update.callback_query
            user_id = update.callback_query['from']['id']
        else:
            return
        lang = UserPoints.get_language_by_user_id(user_id)
        i18n = I18nHelper(lang)
        keyboard_factory = KeyboardFactory(i18n)
        reply_markup = keyboard_factory.create_keyboard("menu")
        dialog = i18n.get_dialog('menu')
        await target.message.reply_text(text = escape(dialog), parse_mode="MarkdownV2", reply_markup=reply_markup)
    except Exception as e:
        logging.error(f"Error in message_menu: {e}")
        collect_error(e)
        return False

async def message_unveil_or_not(update, token, role: str):
    try:
        if update.message:
            target = update
            user_id = update.message['from']['id']
        elif update.callback_query:
            target = update.callback_query
            user_id = update.callback_query['from']['id']
        else:
            return
        lang = UserPoints.get_language_by_user_id(user_id)
        i18n = I18nHelper(lang)
        keyboard_factory = KeyboardFactory(i18n)
        reply_markup = keyboard_factory.create_keyboard("unveil_or_not", token=token, role=role)
        dialog = i18n.get_dialog('unveil_or_not')
        await target.message.reply_text(text = escape(dialog), parse_mode="MarkdownV2", reply_markup=reply_markup)
    except Exception as e:
        logging.error(f"Error in message_unveil_or_not: {e}")
        collect_error(e)
        await telegram_app.bot.send_message(
            chat_id=user_id,
            text='Sorry, I am not able to generate content for you right now. Please try again later.'
        )

        return False
    
def validate_ticker(ticker: str) -> bool:
    """
    验证加密货币缩写是否有效。
    
    参数:
    - ticker: 用户输入的加密货币缩写
    
    返回:
    - 如果有效返回 True，否则返回 False
    """
    # 去除前后的空白字符
    ticker = ticker.strip()
    
    # 检查输入长度
    if len(ticker) < 1 or len(ticker) > 10:  # 假设合理的长度范围是1到10个字符
        return False
    
    # 检查是否为字母和数字的组合
    if not re.match(r'^[A-Za-z0-9]+$', ticker):
        return False
    
    return True