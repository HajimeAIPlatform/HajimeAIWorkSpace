import logging
import asyncio
from io import BytesIO
from os import getenv
import urllib.parse
import time
import json
import os
from typing import List, Dict, Union

from flask import Blueprint, jsonify, request
from md2tgmd import escape
from telegram import Update, BotCommand, InlineKeyboardButton, InlineKeyboardMarkup, WebAppInfo, \
    InlineQueryResultsButton, InputMediaPhoto, Message, CallbackQuery
import PIL as Pillow
from telegram.ext import ApplicationBuilder, DictPersistence, CommandHandler

from pythonp.apps.tokenfate.src.dify.views import chat_blocking, chat_streaming, chat_workflow
from pythonp.apps.tokenfate.src.binance.views import handle_binance_command
from pythonp.apps.tokenfate.src.binance.utils import get_all_prices, process_recommendation
import pythonp.apps.tokenfate.src.ton.views as ton_module
from pythonp.apps.tokenfate.src.bot.commands import set_bot_commands_handler
from pythonp.apps.tokenfate.src.bot.wallet_menu_callback import set_handlers
import pythonp.apps.tokenfate.src.bot.state as ChatStatus
from pythonp.apps.tokenfate.src.bot.i18n_helper import I18nHelper
from pythonp.apps.tokenfate.src.bot.keyboards import KeyboardFactory
from pythonp.apps.tokenfate.src.ton.tc_storage import DailyFortune, UserActivityTracker
from pythonp.apps.tokenfate.models.transaction import UserPoints


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

# 初始化语言和键盘工厂
i18n = I18nHelper()
keyboard_factory = KeyboardFactory(i18n)

level_photo = {
    '上签': '决而能和.png',
    '中签': '应时而变.png',
    '下签': '蓄养待进.png',
}

WEB_MINI_APP_URL = getenv('WEB_MINI_APP_URL')


async def run_bot():
    await set_bot_commands_handler(telegram_app)
    await telegram_app.initialize()
    await telegram_app.start()


asyncio.run(run_bot())

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

    except Exception as error:
        logging.error(f"Error occurred during streaming: {error}")
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

    except Exception as error:
        logging.error(f"Error occurred during streaming: {error}")
        return {
            "text":
                'Sorry, I am not able to generate content for you right now. Please try again later.'
        }


@bot.route('/binance/price', methods=['GET'])
async def binance_price():
    return get_all_prices(10000, "USDT")


async def start(update):
    # 创建没有参数的按钮
    web_app_info = WebAppInfo(url=WEB_MINI_APP_URL + "/#/", api_kwargs={"aa": "11"})
    button = InlineQueryResultsButton(text="Open Wallet", web_app=web_app_info)

    # 创建带有参数的按钮
    action = "buy"
    token = "ORN"
    amount = 100
    timestamp = int(time.time())  # 添加时间戳
    web_app_url_with_params = f"{WEB_MINI_APP_URL}/#/transaction?action={urllib.parse.quote(action)}&token={urllib.parse.quote(token)}&amount={urllib.parse.quote(str(amount))}&t={timestamp}"
    web_app_info2 = WebAppInfo(url=web_app_url_with_params, api_kwargs={"aa": "22"})
    button2 = InlineQueryResultsButton(text="111 Wallet", web_app=web_app_info2)

    keyboard = InlineKeyboardMarkup([[button], [button2]])

    await update.message.reply_text("Click a button to open Web App", reply_markup=keyboard)


@bot.route('/webhook', methods=['POST'])
# @run_async
async def webhook():
    chat_id = None
    try:
        body = request.get_json()
        update = Update.de_json(body, telegram_app.bot)
        logging.debug(f"Received update: {update}")
        await handle_daily_checkin(update)
        if update.edited_message:
            return jsonify({'status': 'ok'}), 200

        if update.callback_query:
            await set_handlers(update, telegram_app)
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
                    parse_mode='HTML'
                )
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data == "show_aura_rules":
                await update.callback_query.answer()
                await show_aura_rules(update.callback_query)
                return jsonify({'status': 'ok'}), 200

            if update.callback_query.data == "connect_wallet_button":
                await update.callback_query.answer()
                # connect to wallet
                user_id = update.callback_query['from']['id']
                connected_wallets = await user_activity_tracker.get_connected_wallets(user_id)
                
                if connected_wallets:
                    # Format the set into a list with each wallet on a new line
                    wallet_list = "\n".join(connected_wallets)
                    await update.callback_query.message.reply_text(f"已连接过的钱包:\n{wallet_list}")
                await ton_module.wallet_menu_callback.on_choose_wallet_click(update)
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data == "aura_action_daily_checkin":
                await update.callback_query.answer()
                user_id = update.callback_query['from']['id']
                if await user_activity_tracker.is_checked_in_today(user_id=user_id):
                    await update.callback_query.message.reply_text("灵气无波，重访无效")
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data == "aura_action_recommend_click":
                await update.callback_query.answer()
                reply_markup = keyboard_factory.create_keyboard("launch")
                await update.callback_query.message.reply_text("完成一次符命探寻，方可征得启示", reply_markup=reply_markup)
                return jsonify({'status': 'ok'}), 200

            if update.callback_query.data.startswith("reveal_fate"):
                data = update.callback_query.data
                logging.info(f"data: {data}")
                details = data.split(":")
                token = details[1]
                token_from = details[2]
                await risk_preference(update, token, token_from)
                return jsonify({'status': 'ok'}), 200
            
            if update.callback_query.data.startswith("risk"):
                data = update.callback_query.data
                logging.info(f"data: {data}")
                details = data.split(":")
                token = details[1]
                token_from = details[2]
                await reveal_fate(update, token, token_from)
                return jsonify({'status': 'ok'}), 200
            

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
            reply_markup = keyboard_factory.create_keyboard("start")
            dialog = i18n.get_dialog('start')
            image_path = get_image_path('Welcome.png')
            with open(image_path, 'rb') as image_file:
                await update.message.reply_photo(
                    photo=image_file, 
                    caption=escape(dialog), 
                    parse_mode="MarkdownV2", 
                    reply_markup=reply_markup
                )  
            return jsonify({'status': 'ok'}), 200

        if update.message.text == '/quest':
            dialog = i18n.get_dialog('quest')
            await update.message.reply_text(
                text=dialog,
                parse_mode='HTML'
            )
            return jsonify({'status': 'ok'}), 200
        
        if update.message.text == '/aura':
            await show_aura_rules(update)
            return jsonify({'status': 'ok'}), 200

        if update.message.text.startswith('$'):
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
                await update.message.reply_text(escape(dialog), parse_mode="HTML", reply_markup=reply_markup)
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
            image = Pillow.Image.open(bytesIO)
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

    except Exception as error:
        logging.error(f"Error Occurred: {error}")
        return {
            "method":
                "sendMessage",
            "chat_id":
                chat_id,
            "text":
                'Sorry, I am not able to generate content for you right now. Please try again later.'
        }


def get_image_path(image_name):
    # 获取项目根目录
    project_root = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
    # 构建图片的绝对路径
    image_path = os.path.join(project_root, 'static', 'images', image_name)
    return image_path

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


def fetch_trending_tokens():
    try:
        query = ""
        
        data = {
            "query": query,
        }

        # Fetch recommendations from Dify API
        api_response = chat_workflow(data)
        logging.info("Raw API response: %s", api_response)
        parsed_response = parse_token_response(api_response)
        logging.info(f"parsed_response: {parsed_response}")
        return validate_token_data(parsed_response)

    except Exception as e:
        logging.error("Failed to fetch trending tokens: %s", e)
        return []


def create_token_keyboard(tokens):
    keyboard = []
    row = []

    for token_info in tokens:
        token_name = token_info.get('token', '')
        if token_name:
            button = InlineKeyboardButton(
                text=f"${token_name}",
                callback_data=f"reveal_fate:${token_name}:recommended"
            )
            row.append(button)

            if len(row) == 5:  # 5 tokens per row
                keyboard.append(row)
                row = []

    if row:  # Handle remaining buttons
        keyboard.append(row)

    return InlineKeyboardMarkup(keyboard)


async def reveal_fate(update, token, token_from='normal'):
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
        if not token:
            await target.message.reply_text("请输入你的符币代码。")
            return jsonify({'status': 'ok'}), 200
        
        # 随机抽签
        chat_id = target.message.chat_id
        cached_lot = await daily_fortune.get_cached_lot(chat_id, token)
        is_cached = cached_lot is not None
        if is_cached:
            result_of_draw = cached_lot
        else:
            result_of_draw = await daily_fortune.get_daily_lot(chat_id, token)
        logging.info(f"result_of_draw: {result_of_draw}")

        if is_cached:
            # 用户积分无需更新
            await get_aura_status(update, "aura_action_cached_hit")
        else:
            # 更新用户积分
            if token_from == 'recommended':
                sql_status = UserPoints.update_points_by_user_id(user_id=user_id, points=10)
                if not sql_status:
                    logging.error("Failed to update user points")
                    await get_aura_status(update, "aura_action_invalid")
                await get_aura_status(update, "aura_action_recommend_click")

            sql_status = UserPoints.update_points_by_user_id(user_id=user_id, points=-15)
            if not sql_status:
                logging.error("Failed to update user points")
                await get_aura_status(update, "aura_action_invalid")
            await get_aura_status(update, "aura_action_fate_reveal")

        sign_level = result_of_draw["sign_level"]
        sign_from = result_of_draw["sign_from"]
        sign_text = result_of_draw["sign_text"]
        # 发送抽签结果
        words = f"**{token}的今日运势**\n\n{sign_text}\n\n摘自：《{sign_from}》\n"
        escaped_words = escape(words)
        image_path = get_image_path(level_photo[sign_level])
        with open(image_path, 'rb') as image_file:
            await target.message.reply_photo(
                photo=image_file,
                caption=escaped_words,
                parse_mode="MarkdownV2",
            )

        # 新的键盘布局：推荐tokens
        recommended_tokens = fetch_trending_tokens()
        if not recommended_tokens:
            logging.error("Failed to get token recommendations")
            return jsonify({'status': 'ok'}), 200
        reply_markup = create_token_keyboard(recommended_tokens)

        # 发送推荐的token
        dialog = i18n.get_dialog("recommended")
        await target.message.reply_text(
            escape(dialog), 
            parse_mode="MarkdownV2", 
            reply_markup=reply_markup
        )

        return jsonify({'status': 'ok'}), 200
    
    except Exception as e:
        logging.error(f"Error in reveal_fate: {e}")
        await target.message.reply_text(
            "Sorry, something went wrong while processing your request."
        )
        return
    
async def risk_preference(update, token, token_from='normal'):
    try:
        await update.callback_query.answer()
        if not token:
            await update.callback_query.message.reply_text("Please Enter Your Token.")
            return jsonify({'status': 'ok'}), 200
        
        # 发送文本信息
        dialog = i18n.get_dialog('risk')
        await update.callback_query.message.reply_text(
            text=dialog,
            parse_mode="MarkdownV2"
        )

        # 发送图片及选择项
        reply_markup = keyboard_factory.create_keyboard("risk", token=token, token_from=token_from)
        image_path = get_image_path('risk_preference_combined.png')
        with open(image_path, 'rb') as image_file:
            await update.callback_query.message.reply_photo(
                photo=image_file,
                parse_mode="MarkdownV2",
                reply_markup=reply_markup
            )
        return jsonify({'status': 'ok'}), 200

    except Exception as e:
        logging.error(f"Error in reveal_fate: {e}")
        await update.callback_query.message.reply_text(
            "Sorry, something went wrong while processing your request."
        )
        return
    
async def show_aura_rules(update: Union[Update, CallbackQuery]):
    try:
        if isinstance(update, CallbackQuery):
            await update.answer()
            target = update.message
            user_id = update['from']['id']
        elif isinstance(update, Update):
            target = update.message
            user_id = update.message['from']['id']
        else:
            return

        # Send aura rules information
        dialog = i18n.get_dialog('aura_rules')
        dialog = dialog.format(points= UserPoints.get_points_by_user_id(user_id=user_id))
        image_path = get_image_path("吾之灵气.png")
        with open(image_path, 'rb') as image_file:
            await target.reply_photo(
                photo=image_file,
                caption=escape(dialog),
                parse_mode="MarkdownV2",
            )

        # Send interactive elements
        reply_markup = keyboard_factory.create_keyboard("aura")
        image_path = get_image_path('ways-to-impact-aura.png')
        with open(image_path, 'rb') as image_file:
            await target.reply_photo(
                photo=image_file,
                reply_markup=reply_markup
            )
        return jsonify({'status': 'ok'}), 200

    except Exception as e:
        logging.error(f"Error in show_aura_rules: {e}")
        if isinstance(update, CallbackQuery):
            await update.message.reply_text("Sorry, something went wrong while processing your request.")
        elif isinstance(update, Message):
            await update.reply_text("Sorry, something went wrong while processing your request.")
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
        dialog = i18n.get_dialog(action)
        dialog = dialog.format(user_id=user_id)
        button = i18n.get_button("aura")
        reply_markup = InlineKeyboardMarkup([[InlineKeyboardButton(text=button, callback_data="show_aura_rules")]])
        await target.message.reply_text(escape(dialog), parse_mode="MarkdownV2", reply_markup=reply_markup)
        return jsonify({'status': 'ok'}), 200

    except Exception as e:
        logging.error(f"Error in get_aura_status: {e}")
        await target.message.reply_text("Sorry, something went wrong while processing your request.")
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
        else:
            logging.error("Failed to check in user")

    except Exception as e:
        logging.error(f"Error in handle_daily_checkin: {str(e)}")
        await get_aura_status(update, "aura_action_invalid")
