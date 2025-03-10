import queue
import threading
import logging
import json
import time
from flask import current_app
from pythonp.apps.fortune_teller.models.transaction import save_binance_transaction_to_db
from pythonp.apps.fortune_teller.service.binance.utils import (get_binance_client,
                               get_common_currency_price, create_order,
                               get_order_info, get_result_info)

# 初始化队列
transaction_queue = queue.Queue()


def process_transaction_queue(app):
    while True:
        side, symbol, amount, transaction_pair_id = transaction_queue.get()
        try:
            # trade logic
            result = create_order(symbol, side, amount)
            # 提取所需的值
            symbol, order_id, status, type, side, transact_time, cummulative_quote_qty, total_commission, commission_details = get_result_info(
                result)
            fee = total_commission

            with app.app_context():
                save_binance_transaction_to_db(side=side,
                                               status=status,
                                               symbol=symbol,
                                               amount=amount,
                                               cummulative_quote_qty=cummulative_quote_qty,
                                               fee=fee,
                                               type=type,
                                               transaction_pair_id=transaction_pair_id,
                                               order_id=order_id,
                                               timestamp=transact_time,
                                               full_data=json.dumps(result))

        except Exception as e:
            logging.error(f"Error processing binance transaction: {e}")
            with app.app_context():
                transact_time = int(time.time() * 1000)
                save_binance_transaction_to_db(side=side,
                                               status="ERROR",
                                               symbol=symbol,
                                               amount=amount,
                                               cummulative_quote_qty=0.0,
                                               fee=0.0,
                                               type="ERROR",
                                               transaction_pair_id=transaction_pair_id,
                                               order_id="-1",
                                               timestamp=transact_time,
                                               full_data=json.dumps("trade error, need to handle"))
        finally:
            transaction_queue.task_done()



def put_transaction_into_queue(side, symbol, amount, transaction_pair_id):
    logging.info(f"Putting transaction {transaction_pair_id} into queue")
    transaction_queue.put((side, symbol, amount, transaction_pair_id))


def start_transaction_processor(app):
    transaction_processor_thread = threading.Thread(
        target=process_transaction_queue,args=(app,), daemon=True)
    transaction_processor_thread.start()
