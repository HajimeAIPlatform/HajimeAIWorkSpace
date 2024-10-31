import os
from telegram import InlineKeyboardButton, InlineKeyboardMarkup, InputFile, Update, InputMediaPhoto, InputMedia
from telegram.ext import ApplicationBuilder, CallbackQueryHandler, CommandHandler, ContextTypes
import qrcode
from io import BytesIO
import asyncio
from urllib.parse import urlencode
from src.ton.connector import get_connector
import logging
from pytoniq_core import Address
from pytonconnect import TonConnect
import src.bot.state as TaskState
import src.ton.views as ton_module

def get_wallets():
    wallets_list = TonConnect.get_wallets()
    return wallets_list


async def send_callback(chat_id, connector, telegram_app):
    try:
        await wait_for_connection(connector, telegram_app, chat_id)
    except Exception as e:
        logging.error(f"Error in send_callback: {e}")
    finally:
        TaskState.remove_waiting_task(chat_id)


async def wait_for_connection(connector, telegram_app, chat_id, timeout=180):
    TaskState.add_waiting_task(chat_id, True)
    for i in range(1, timeout):
        if TaskState.get_waiting_task(chat_id) is None:
            return
        await asyncio.sleep(1)

        # if i == timeout/2:
        print(
            f'Waiting for connection... {i}, connector.connected : {connector.connected}'
        )
        if connector.connected:
            if connector.account.address:
                wallet_address = connector.account.address
                wallet_address = Address(wallet_address).to_str(
                    is_bounceable=False)
                logging.info(f'Connected with address: {wallet_address}')
                if wallet_address:
                    await telegram_app.bot.send_message(
                        chat_id=chat_id,
                        text=
                        f'You are connected with address <code>{wallet_address}</code>',
                        parse_mode='HTML')
                    query_token = TaskState.get_tmp_token(chat_id)
                    TaskState.remove_tmp_token(chat_id)
                    if query_token:
                        # 构造新的文本内容
                        escaped_text = f'Thank you, my dear.\nNow, with the stars aligned, it is time to reveal the fate of your token.'

                        # 重新创建原有的两个按钮
                        button = InlineKeyboardButton(text="Reveal My Token Fate Now", callback_data=f"reveal_fate:{query_token}")

                        # 创建新的键盘布局
                        new_keyboard = InlineKeyboardMarkup([[button]])

                        # 根据您的模板发送消息
                        await telegram_app.bot.send_message(
                            chat_id=chat_id,
                            text=escaped_text,
                            parse_mode='HTML',
                            reply_markup=new_keyboard
                        )
                    return 
    await telegram_app.bot.send_message(chat_id=chat_id,
                                        text='Connect wallet timeout!')


ITEMS_PER_PAGE = 3


async def on_choose_wallet_click(update: Update):
    logging.info("Starting on_choose_wallet_click")
    if update.callback_query is not None:
        query = update.callback_query
        data = query.data
        chat_id = query.message.chat_id
    else:
        # 处理没有 callback_query 的情况，例如通过命令触发
        query = None
        data = update.message.text
        chat_id = update.message.chat_id

    TaskState.remove_waiting_task(chat_id)

    # 获取当前页码，默认为 0
    if ':' in data:
        current_page = int(data.split(':')[1])
    else:
        current_page = 0

    wallets = get_wallets()

    # 计算总页数
    total_pages = (len(wallets) + ITEMS_PER_PAGE - 1) // ITEMS_PER_PAGE

    # 获取当前页的钱包
    start_idx = current_page * ITEMS_PER_PAGE
    end_idx = start_idx + ITEMS_PER_PAGE
    page_wallets = wallets[start_idx:end_idx]

    keyboard = [[
        InlineKeyboardButton(
            wallet['name'],
            callback_data=f"select_wallet:{wallet['app_name']}")
        for wallet in page_wallets
    ]]

    # 添加分页按钮
    pagination_buttons = []
    if current_page > 0:
        pagination_buttons.append(InlineKeyboardButton("« Left", callback_data=f"choose_wallet:{current_page - 1}"))
    if current_page < total_pages - 1:
        pagination_buttons.append(InlineKeyboardButton("Right »", callback_data=f"choose_wallet:{current_page + 1}"))

    if pagination_buttons:
        keyboard.append(pagination_buttons)

    # 添加返回按钮
    # keyboard.append([InlineKeyboardButton("« Back", callback_data="universal_qr")])

    reply_markup = InlineKeyboardMarkup(keyboard)

    # 如果是通过按钮点击触发的，更新消息的按钮
    if query is not None:
        logging.info("Before connect(callback query)")
        # await query.edit_message_caption(caption="Choose a wallet:", reply_markup=reply_markup)
        await callback_query_connect(update, reply_markup)
        logging.info("After connect(callback query)")
        await query.answer()
        logging.info("After query.answer")
    else:
        # 通过命令触发的情况，发送新的消息
        logging.info("Before connect(Command triggered)")
        await connect(update, reply_markup)
        logging.info("After connect(Command triggered)")
    logging.info("Finishing on_choose_wallet_click")
    return True


async def on_open_universal_qr_click(update: Update, telegram_app):
    query = update.callback_query
    chat_id = update.update_id

    wallets = get_wallets()

    connector = get_connector(chat_id)
    link = await connector.connect(wallets)

    await edit_qr(query.message, link, telegram_app, 'Select a wallet')

    keyboard = build_universal_keyboard(link)
    reply_markup = InlineKeyboardMarkup(keyboard)

    await query.edit_message_reply_markup(reply_markup=reply_markup)
    await query.answer()


async def on_wallet_click(update: Update, telegram_app):
    try:
        query = update.callback_query
        data = query.data.split(":")[1] # wallet app_name
        chat_id = query.message.chat_id

        connector = get_connector(chat_id)
        wallets = get_wallets()
        selected_wallet = await get_wallet_info(data)
        print(selected_wallet, 'selected_wallet')
        if not selected_wallet:
            print(wallets, 'wallets')
            return

        button_link = await connector.connect({
            'bridge_url':
                selected_wallet['bridge_url'],
            'universal_url':
                selected_wallet['universal_url']
        })

        qr_link = button_link

        await edit_qr(query.message,
                      qr_link,
                      telegram_app,
                      caption=f"Connect wallet within 3 minutes")

        keyboard = [[
            InlineKeyboardButton("« Back", callback_data="choose_wallet"),
            InlineKeyboardButton(f"Open {selected_wallet['name']}",
                                 url=button_link)
        ]]
        reply_markup = InlineKeyboardMarkup(keyboard)

        await query.edit_message_reply_markup(reply_markup=reply_markup)
        await query.answer()

        asyncio.create_task(send_callback(chat_id, connector, telegram_app))
    except Exception as e:
        logging.error(f"Error in on_wallet_click: {e}")
        await query.answer("Something went wrong. Please try again later.")


async def edit_qr(message, link, telegram_app, caption='', box_size=4, border=4):
    qr_img = qrcode.make(link, box_size=box_size, border=border)
    stream = BytesIO()
    qr_img.save(stream)
    stream.seek(0)  # 将指针移动到开始位置

    # 设置 .name 属性
    stream.name = 'qrcode.png'

    # 现在可以将 stream 作为 media 参数传递
    media = stream

    await telegram_app.bot.edit_message_media(
        chat_id=message.chat_id,
        message_id=message.message_id,
        media=InputMediaPhoto(media=media, caption=caption))


def convert_deeplink_to_universal_link(deeplink, universal_link):
    # Implement the logic to convert deeplink to universal link if necessary
    return universal_link


def add_tg_return_strategy(link, strategy):
    # Implement your logic to add Telegram return strategy to the link
    return f"{link}?return_strategy={strategy}"


def build_universal_keyboard(link):
    wallets = get_wallets()
    keyboard = [
        InlineKeyboardButton('Choose a Wallet', callback_data='choose_wallet'),
        # InlineKeyboardButton('Open Link', url=f'https://ton-connect.github.io/open-tc?connect={urlencode({"connect": link})}')
    ]
    return [keyboard]

async def callback_query_connect(update: Update, reply_markup):
    chat_id = update.callback_query.message.chat_id

    connector = get_connector(chat_id)
    wallets = get_wallets()

    link = await connector.connect(wallets)

    qr_img = qrcode.make(link, box_size=4, border=4)
    stream = BytesIO()
    qr_img.save(stream)
    file = InputFile(stream.getvalue(), filename='qrcode.png')

    await update.callback_query.message.reply_photo(photo=file,
                                     caption='Choose a wallet:',
                                     reply_markup=reply_markup)
    return True

async def connect(update: Update, reply_markup):
    chat_id = update.message.chat_id

    connector = get_connector(chat_id)
    wallets = get_wallets()

    link = await connector.connect(wallets)

    qr_img = qrcode.make(link, box_size=4, border=4)
    stream = BytesIO()
    qr_img.save(stream)
    file = InputFile(stream.getvalue(), filename='qrcode.png')

    await update.message.reply_photo(photo=file,
                                     caption='Choose a wallet:',
                                     reply_markup=reply_markup)
    return True


async def get_wallet_info(wallet_name):
    wallets = get_wallets()
    for wallet in wallets:
        if wallet['app_name'] == wallet_name:
            return wallet


async def handle_send_transaction(update: Update, telegram_app, text, side):
    query = update.callback_query
    chat_id = query.message.chat.id
    if ':' in text:
        currency, amount = text.replace(f'{side.lower()} ', '').split(':')
        print(currency, amount, 'currency, amount')
        try:
            amount = float(amount)
            message = await ton_module.send_transaction(update, telegram_app,
                                                        currency, amount, side)
            await telegram_app.bot.send_message(chat_id=chat_id,
                                                text=message)
        except ValueError:
            await telegram_app.bot.send_message(chat_id=chat_id,
                                                text="Invalid amount. Please enter in the format BTC:0.1,\ncancel to "
                                                     "input /cancel")
    else:
        await telegram_app.bot.send_message(chat_id=chat_id,
                                            text="Invalid format. Please enter in the format BTC:0.1,\n Cancel to "
                                                 "input /cancel.")


async def set_handlers(update, telegram_app):
    query = update.callback_query
    if query.data.startswith('choose_wallet'):
        await on_choose_wallet_click(update)
    elif query.data == 'universal_qr':
        await on_open_universal_qr_click(update, telegram_app)
    elif query.data.startswith('select_wallet'):
        await on_wallet_click(update, telegram_app)
    elif query.data.startswith('buy'):
        await handle_send_transaction(update, telegram_app, query.data, "BUY")
    elif query.data.startswith('sell'):
        await handle_send_transaction(update, telegram_app, query.data, "SELL")
