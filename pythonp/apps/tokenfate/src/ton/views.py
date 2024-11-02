# ton_module.py
import time
import asyncio
import logging
from pytoniq_core import Cell
import pytonconnect.exceptions
from telegram import Update
from telegram.ext import ConversationHandler
from src.ton.messages import get_comment_message
from src.ton.connector import get_connector
import src.bot.wallet_menu_callback as wallet_menu_callback
from models.transaction import save_ton_transaction_to_db, UserAsset
from src.ton.utils import calculate_fee, get_bot_ton_address, convert_address_to_hex, estimate_gas_fee
import src.bot.state as ChatStatus
from src.binance.views import convert_currency
from src.binance.transaction_queue import put_transaction_into_queue
from src.binance.utils import is_min_trade_quantity_limit
import src.ton.ton_server as ton_server
import mock_service.mock_connector as mock_connector
import mock_service.mock_ton_server as mock_ton_server
from os import getenv

MOCK_SERVER = getenv('MOCK_SERVER', 'False').lower() in ('true', '1', 't')

async def disconnect_wallet(update):
    chat_id = update.message.chat_id
    connector = get_connector(chat_id)
    await connector.restore_connection()
    await connector.disconnect()


def save_to_transaction_info_to_db(user_id, chat_id, symbol, side, amount, fee,
                                   deposit_address, trace_link):
    try:
        transaction_pair_id = save_ton_transaction_to_db(
            user_id=user_id,
            chat_id=chat_id,
            symbol=symbol,
            side=side,
            status='success',
            address=deposit_address,
            amount=amount,
            trace_link=trace_link,  # 从上下文中获取
            fee=fee)

        put_transaction_into_queue(side, symbol, amount, transaction_pair_id)
    except Exception as e:
        logging.error(f'Error saving transaction to DB: {e}')


async def send_and_confirm_transaction(user_id, chat_id, connector, telegram_app, symbol,
                                       amount, side, amount_after_fee, fee):
    text = ""
    try:
        deposit_address = await ton_server.get_ton_wallet_address()
        transaction = {
            'valid_until':
                int(time.time() + 3600),
            'messages': [
                get_comment_message(
                    destination_address=deposit_address,
                    amount=amount_after_fee * 10 ** 9,
                    comment=
                    f'Symbol type: {symbol}, Amount: {amount}, Fee: {fee}')
            ]
        }
        logging.info(
            f'User buy {amount} {symbol}.Sending transaction to {deposit_address} with amount: {amount_after_fee}'
        )
        print(f"MOCK_SERVER: {MOCK_SERVER}, Type: {type(MOCK_SERVER)}")

        server = mock_connector if MOCK_SERVER else connector
        result = await asyncio.wait_for(
            server.send_transaction(transaction=transaction), 300)
        print(result, 'result')
        print(result["boc"])
        text = 'Transaction sent successfully!'
        msg_hash = Cell.one_from_boc(result["boc"]).hash.hex()
        trace_link = f"https://toncenter.com/api/v3/transactionsByMessage?direction=out&msg_hash={msg_hash}&limit=128&offset=0"
        print(f"Transaction info -> {trace_link}")
        # 发送交易后保存到数据库
        save_to_transaction_info_to_db(user_id, chat_id,
                                       symbol, side, amount, fee,
                                       deposit_address, trace_link)

    except asyncio.TimeoutError:
        text = 'Transaction Timeout error!'
    except pytonconnect.exceptions.UserRejectsError:
        text = 'You rejected the transaction!'
    except Exception as e:
        text = f'Transaction Unknown error: {e}'

    await telegram_app.bot.send_message(chat_id=chat_id, text=text)
    ChatStatus.remove_transaction_status(chat_id)


async def transfer_ton_to_user(user_id, chat_id, telegram_app, symbol,
                               amount_ton, amount, side, fee,
                               destination_address):
    text = ""
    try:
        total_gas_fee_in_ton = estimate_gas_fee()

        amount_after_fee = float(amount_ton) - float(fee) - float(
            total_gas_fee_in_ton)
        amount_after_fee = round(amount_after_fee, 9)
        server = mock_ton_server if MOCK_SERVER else ton_server
        tx_hash = await server.send_transaction(
            destination_address, amount_after_fee,
            f"{side} {symbol} cost, amount:{amount}")

        if tx_hash:
            logging.info(f"Transaction {tx_hash} sent!")
            # Format the message
            message = (f"Your sell assets have arrived\n"
                       f"Sell {symbol.upper()}: {amount}\n"
                       f"Received TON: {amount_after_fee}\n"
                       f"Transaction hash: `{tx_hash}`")
            save_to_transaction_info_to_db(user_id,
                                           chat_id, symbol, side, amount, fee,
                                           destination_address, tx_hash)
            logging.info(message)
            # Send the message
            await telegram_app.bot.send_message(chat_id=chat_id,
                                                text=message,
                                                parse_mode='Markdown')

        else:
            logging.error(f"Failed to send transaction")
            await telegram_app.bot.send_message(
                chat_id=chat_id, text="Failed to send transaction")
    except Exception as e:
        logging.error(f"Error in transfer_ton_to_user: {e}")
        await telegram_app.bot.send_message(chat_id=chat_id,
                                            text="Error in sell assets")
    finally:
        ChatStatus.remove_transaction_status(chat_id)


async def send_transaction(update, telegram_app, symbol, amount, side):
    try:
        if update.callback_query:
            query = update.callback_query
            chat_id = query.message.chat_id
            user_id = query.from_user.id
        else:
            chat_id = update.message.chat_id
            user_id = update.message.from_user.id

        connector = get_connector(chat_id)
        connected = await connector.restore_connection()
        if not connected:
            return 'Connect wallet first!'
        is_true, minimum = is_min_trade_quantity_limit(symbol, amount)
        if is_true:
            return f'The minimum trade quantity is {minimum} {symbol}'
        #     await update.message.reply_text(f'The minimum trade quantity is {minimum} {symbol}')
        # return WAITING_FOR_INPUT

        amount_ton,from_price, to_price = convert_currency(symbol, "TON", amount)
        fee = await calculate_fee(float(amount_ton))
        print(f"Total fee: {fee} TON,  {symbol} to TON, from_amount: {amount},to_amount: {amount_ton}")
        # buy
        if side == 'BUY':
            amount_after_fee = float(amount_ton) + float(fee)
            asyncio.create_task(
                send_and_confirm_transaction(user_id, chat_id, connector, telegram_app,
                                             symbol, amount, side,
                                             amount_after_fee, fee))
            return (f'Sending transaction...\n'
                    f'Approve transaction in your wallet app!\n\n'
                    f'Transaction Details:\n'
                    f'{side.upper()} {symbol} : {amount}\n'
                    f'Price : {from_price}')

        elif side == 'SELL':
            # search assets
            asset_amount = UserAsset.get_amount_by_user_id_and_symbol(
                update.message.from_user.id, symbol)
            if asset_amount is None or asset_amount < float(amount):
                return 'You do not have enough assets!'

            asyncio.create_task(
                transfer_ton_to_user(user_id, chat_id, telegram_app, symbol,
                                     amount_ton, amount, side, fee,
                                     connector.account.address))
            return (f'Confirming the transaction, please look for a reply later or check your wallet assets!\n\n'
                    f'Transaction Details:\n'
                    f'{side.upper()} {symbol} : {amount}\n'
                    f'Price : {from_price}')
        else:
            return 'Invalid side!'
    except Exception as e:
        logging.error(f"Error in send_transaction: {e}")
        return f"{e}"


async def check_connected(update, telegram_app):
    if update.message:
        chat_id = update.message.chat_id
    elif update.callback_query:
        chat_id = update.callback_query.message.chat_id
    logging.info(f"Checking wallet connection for chat_id: {chat_id}")

    # Get connector instance
    connector = get_connector(chat_id)
    logging.info(f"Connector: {connector}")
    if not connector:
        await telegram_app.bot.send_message(
            chat_id=chat_id,
            text='Unable to initialize wallet connector. Please try again later.',
            parse_mode='HTML'
        )
        return None
    connected = None
    # Attempt to restore connection
    try:
        connected = await connector.restore_connection()
    except Exception as e:
        logging.error(f"Error restoring connection: {e}")
    logging.info(f"connected == {connected}")
    if not connected:
        return None

    # if connector.account and connector.account.address:
    #     wallet_address = convert_address_to_hex(connector.account.address)
    #     await telegram_app.bot.send_message(
    #         chat_id=chat_id,
    #         text=
    #         f'You are connected with address <code>{wallet_address}</code>',
    #         parse_mode='HTML')
    return True


async def get_my_wallet(update, telegram_app):
    chat_id = update.message.chat_id
    connector = get_connector(chat_id)
    connected = await connector.restore_connection()
    if not connected:
        return {"method": "sendMessage", "text": 'Connect wallet first!'}

    if connector.account and connector.account.address:
        wallet_address = convert_address_to_hex(connector.account.address)
        assets = UserAsset.get_assets_by_user_id(update.message.from_user.id)

        if assets:
            header = "Your Assets:\n"
            table_rows = "\n".join(f"{symbol:<10} {amount}" for symbol, amount in assets.items())
            table_message = (
                    f'You are connected with address `{wallet_address}`\n\n' +
                    header + table_rows
            )
        else:
            table_message = f'You are connected with address `{wallet_address}`'

        await telegram_app.bot.send_message(
            chat_id=chat_id,
            text=table_message,
            parse_mode='Markdown'
        )

        return True
    else:
        return {"method": "sendMessage", "text": 'No wallet connected!'}


async def send_tx(update, context):
    chat_id = update.message.chat_id
    connector = get_connector(chat_id)
    connected = await connector.restore_connection()
    if not connected:
        await update.message.reply_text(
            "Connect wallet first!")
        return ConversationHandler.END
    ChatStatus.set_transaction_status(chat_id, 'transition')
    await update.message.reply_text(
        "Please enter the currency and amount in the format BTC:0.1")
    return WAITING_FOR_INPUT


async def handle_ton_command(telegram_app, update):
    # if command.startswith('/connect'):
    if update.message is None:
        return None
    command = update.message.text
    if command == '/connect':
        is_connected = await check_connected(update, telegram_app)
        if is_connected:
            return True
        isTrue = await wallet_menu_callback.on_choose_wallet_click(update)

        return isTrue
    elif command == '/disconnect':
        chat_id = update.message.chat_id
        connector = get_connector(chat_id)
        connected = await connector.restore_connection()
        if not connected:
            return {
                "method": "sendMessage",
                "text": "Connect wallet first!"
            }
        await disconnect_wallet(update)
        return {
            "method": "sendMessage",
            "text": 'You have been successfully disconnected!'
        }

    elif command == '/my_wallet':
        res = await get_my_wallet(update, telegram_app)
        return res
    logging.info(f"Unknown command: {command}")
    return None  # 如果没有匹配的命令，返回 None


WAITING_FOR_INPUT = range(1)  # Define your state


async def handle_send_transaction(update: Update, context, side):
    chat_id = update.message.chat_id
    telegram_app = context.bot
    text = update.message.text
    if ':' in text:
        currency, amount = text.split(':')
        try:
            amount = float(amount)
            message = await send_transaction(update, telegram_app,
                                             currency, amount, side)
            await update.message.reply_text(f"{message}")
        except ValueError:
            await update.message.reply_text(
                "Invalid amount. Please enter in the format BTC:0.1,\ncancel to input /cancel"
            )
            return WAITING_FOR_INPUT
    else:
        await update.message.reply_text(
            "Invalid format. Please enter in the format BTC:0.1,\n Cancel to input /cancel."
        )
        return WAITING_FOR_INPUT
    return ConversationHandler.END


async def buy_transaction(update: Update, context):
    return await handle_send_transaction(update, context, "BUY")


async def sell_transaction(update: Update, context):
    return await handle_send_transaction(update, context, "SELL")


async def cancel(update: Update, context):
    chat_id = update.message.chat_id
    ChatStatus.remove_transaction_status(chat_id)
    await update.message.reply_text("Transaction cancelled.")
    return ConversationHandler.END

async def connect_ton_wallet(chat_id):
    connector = get_connector(chat_id)
    try:
        connected = await connector.restore_connection()
        if connected:
            # You might want to save the connection status or do additional setup here
            return True
        else:
            return False
    except Exception as e:
        logging.error(f"Error connecting TON wallet: {e}")
        return False