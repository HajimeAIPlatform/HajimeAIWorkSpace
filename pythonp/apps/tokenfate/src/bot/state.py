from collections import defaultdict
user_chat_status = {}
waiting_tasks = {}
tmp_tokens = {}
chat_free = {}
user_checkins = {}
user_wallets = defaultdict(set)

def get_transaction_status(chat_id):
    
    return user_chat_status.get(chat_id, None)


def set_transaction_status(chat_id, status):
    user_chat_status[chat_id] = status


def remove_transaction_status(chat_id):
    if chat_id in user_chat_status:
        del user_chat_status[chat_id]

def get_waiting_task(chat_id):
    return waiting_tasks.get(chat_id, None)
def add_waiting_task(chat_id, task):
    waiting_tasks[chat_id] = task

def remove_waiting_task(chat_id):
    if chat_id in waiting_tasks:
        del waiting_tasks[chat_id]

def get_tmp_token(chat_id):
    return tmp_tokens.get(chat_id, None)
def set_tmp_token(chat_id, token):
    tmp_tokens[chat_id] = token  
def remove_tmp_token(chat_id):
    if chat_id in tmp_tokens:
        del tmp_tokens[chat_id]

def get_free(chat_id):
    return chat_free.get(chat_id, False)
def set_free(chat_id):
    chat_free[chat_id] = True 
def remove_free(chat_id):
    if chat_id in chat_free:
        del chat_free[chat_id]        