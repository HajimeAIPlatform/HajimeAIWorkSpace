import logging
import json
import time
from os import getenv

import asyncio
from flask import Blueprint, jsonify, request
from pytoniq_core import Cell
import pytonconnect.exceptions


from pythonp.apps.tokenfate.src.ton.tc_storage import TcStorage, TransactionManager
from pythonp.apps.tokenfate.src.ton.connector import get_connector
from pythonp.apps.tokenfate.src.binance.utils import is_min_trade_quantity_limit
from pythonp.apps.tokenfate.src.binance.views import convert_currency
from pythonp.apps.tokenfate.src.ton.utils import calculate_fee, get_bot_ton_address, convert_address_to_hex, estimate_gas_fee
from pythonp.apps.tokenfate.models.transaction import save_ton_transaction_to_db, UserAsset
import pythonp.apps.tokenfate.src.ton.ton_server as ton_server
from pythonp.apps.tokenfate.src.ton.messages import get_comment_message
import pythonp.apps.tokenfate.mock_service.mock_connector as mock_connector
import pythonp.apps.tokenfate.mock_service.mock_ton_server as mock_ton_server
from pythonp.apps.tokenfate.src.ton.views import save_to_transaction_info_to_db


ton = Blueprint('ton', __name__)

MOCK_SERVER = getenv('MOCK_SERVER', 'False').lower() in ('true', '1', 't')


def get_data(data):
    session_data = {
        "type": data.get("type", "http"),
        "session": {
            "session_private_key": data["session"]["sessionKeyPair"]["secretKey"],
            "wallet_public_key": data["session"]["walletPublicKey"],
            "bridge_url": data["session"]["bridgeUrl"]
        },
        "last_wallet_event_id": data["lastWalletEventId"],
        "connect_event": {
            "id": data["connectEvent"]["id"],
            "event": data["connectEvent"]["event"],
            "payload": {
                "items": [
                    {
                        "name": item["name"],
                        "address": item["address"],
                        "network": item["network"],
                        "publicKey": item["publicKey"],
                        "walletStateInit": item["walletStateInit"]
                    } for item in data["connectEvent"]["payload"]["items"]
                ],
                "device": data["connectEvent"]["payload"]["device"]
            }
        },
        "next_rpc_request_id": data["nextRpcRequestId"]
    }
    return json.dumps(session_data)


async def save_session_data(chat_id, session_data, last_event_id):
    try:
        storage = TcStorage(chat_id)
        await storage.set_item("last_event_id", last_event_id)
        await storage.set_item("connection", session_data)

        connector = get_connector(chat_id)
        connected = await connector.restore_connection()
        if not connected:
            logging.info("Failed to restore connection")
            return None

        return connected
    except Exception as e:
        logging.info(f"Error saving session data: {e}")
        return None


@ton.route('/connect', methods=['POST'])
async def user_connect():
    try:
        data = request.get_json()
        print(data,'data')
        session_data = data.get('session_data')
        last_event_id = data.get('last_event_id')
        chat_id = data.get('chat_id')

        chat_id = int(chat_id)

        session_data_json = get_data(session_data)

        res = await save_session_data(chat_id, session_data_json, last_event_id)
        if res is None:
            return jsonify({"status": "error", "message": "Error to connect wallet"})

        return jsonify({"status": "success", "message": "connect successful"}), 200
    except Exception as e:
        return jsonify({"status": "error", "message": str(e)}), 500


@ton.route('/disconnect', methods=['POST'])
async def user_disconnect():
    try:
        data = request.get_json()
        chat_id = data.get('chat_id')
        chat_id = int(chat_id)


        connector = get_connector(chat_id)
        await connector.disconnect()

        return jsonify({"status": "success", "message": "disconnect successful"}), 200
    except Exception as e:
        return jsonify({"status": "error", "message": str(e)}), 500

@ton.route('/my_wallet', methods=['GET'])
async def user_wallet():
    try:
        user_id = request.args.get('user_id')
        print(user_id,'user_id')
        assets = UserAsset.get_assets_by_user_id(int(user_id))

        return jsonify({"status": "success", "assets": assets,"message":"get successful"}), 200
    except Exception as e:
        return jsonify({"status": "error", "message": str(e)}), 500

async def send_and_confirm_transaction(user_id, chat_id, connector, symbol,
                                       amount, side, amount_after_fee, fee,transaction_code):
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
        print("===========",amount,side,amount_after_fee,fee)
        result = await asyncio.wait_for(
            server.send_transaction(transaction=transaction), 180)

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
        status = "success"
    except asyncio.TimeoutError:
        text = 'Transaction Timeout error!'
        status = "error"
    except pytonconnect.exceptions.UserRejectsError:
        text = 'You rejected the transaction!'
        status = "error"
    except Exception as e:
        text = f'Transaction Unknown error: {e}'
        status = "error"

    await TransactionManager().update_transaction_status(transaction_code, status, text)


async def transfer_ton_to_user(user_id, chat_id, symbol,
                               amount_ton, amount, side, fee,
                               destination_address,transaction_code):
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
            text = "Transaction sent successfully!"
            status = "success"

        else:
            logging.error(f"Failed to send transaction")
            status = "error"
            text = "Failed to send transaction"
    except Exception as e:
        logging.error(f"Error in transfer_ton_to_user: {e}")
        text = "Error in sell assets"
        status="error"
    finally:
        await TransactionManager().update_transaction_status(transaction_code,
                                                            status,
                                                            text)

@ton.route('/transaction', methods=['POST'])
async def user_transaction():
    try:
        data = request.get_json()
        user_id = data.get('user_id')
        chat_id = data.get('chat_id')
        amount = data.get('amount')
        symbol = data.get('symbol')
        side = data.get('side')
        transaction_code = data.get('transaction_code')

        user_id = int(user_id)
        chat_id = int(chat_id)

        connector = get_connector(chat_id)

        connected = await connector.restore_connection()
        if not connected:
            return jsonify({"status": "error", "message": "Connect wallet first!"}), 200

        print(connector.account.address)

        is_true, minimum = is_min_trade_quantity_limit(symbol, amount)
        if is_true:
            return jsonify({"status": "error", "message": f"The minimum trade quantity is {minimum} {symbol}"}), 200

        amount_ton, from_price, to_price = convert_currency(symbol, "TON", amount)
        fee = await calculate_fee(float(amount_ton))
        print(f"Total fee: {fee} TON,  {symbol} to TON, from_amount: {amount},to_amount: {amount_ton}")
        # buy
        side = side.upper()
        if side == 'BUY':
            amount_after_fee = float(amount_ton) + float(fee)
            coro = send_and_confirm_transaction(user_id, chat_id, connector, symbol, amount,
                                                side, amount_after_fee, fee, transaction_code)
            await TransactionManager().start_transaction(transaction_code, coro)

            return jsonify({"status": "success", "message": "transaction initiated"}),200

        elif side == 'SELL':
            # search assets
            asset_amount = UserAsset.get_amount_by_user_id_and_symbol(
                user_id, symbol)
            if asset_amount is None or asset_amount < float(amount):
                return jsonify({"status": "error", "message": "You do not have enough assets!"}),200

            coro = transfer_ton_to_user(user_id, chat_id, symbol, amount_ton, amount, side,
                                        fee, connector.account.address, transaction_code)
            await TransactionManager().start_transaction(transaction_code, coro)
            return jsonify({"status": "success", "message": "Transaction initiated"}), 200

        return jsonify({"status": "success", "message": "transaction successful"}), 200
    except Exception as e:
        return jsonify({"status": "error", "message": str(e)}), 500


@ton.route('/transaction_status', methods=['GET'])
async def check_transaction_status():
    try:
        transaction_code = request.args.get('transaction_code')
        if not transaction_code:
            return jsonify({"status": "error", "message": "Transaction code is required"}), 400

        status = await TransactionManager().get_transaction_status(transaction_code)
        if status:
            return jsonify(status), 200
        else:
            return jsonify({"status": "error", "message": "Transaction not found"}), 404
    except Exception as e:
        return jsonify({"status": "error", "message": str(e)}), 500

@ton.route('/cancel_transaction', methods=['GET'])
async def cancel_transaction():
    try:
        transaction_code = request.args.get('transaction_code')
        cancelled = await TransactionManager().cancel_transaction(transaction_code)
        if cancelled:
            return jsonify({"status": "success", "message": "Transaction cancelled"}), 200
        else:
            return jsonify({"status": "error", "message": "Transaction not found"}), 404
    except Exception as e:
        return jsonify({"status": "error", "message": str(e)}), 500
