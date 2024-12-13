# error_collector.py

import traceback
from pythonp.common.email_notifications.email_global import error_queue

def collect_error(error):
    tb = traceback.extract_tb(error.__traceback__)
    filename, lineno, funcname, _ = tb[-1]
    
    error_info = (
        type(error).__name__,
        str(error),
        filename,
        funcname,
        lineno
    )
    
    error_queue.put(error_info)  # 将错误信息放入队列
    print(f"Added error: {error_info}")  # 确认添加的错误信息
    print_queue_contents()  # 打印当前队列中的所有错误信息

def print_queue_contents():
    # 打印当前队列中的所有错误信息
    print("Current errors in queue:")
    temp_list = []
    print(f"Collecting errors on queue ID: {id(error_queue)}")
    while not error_queue.empty():
        error = error_queue.get()
        temp_list.append(error)  # 将错误信息存储到临时列表中，以便后续使用
        print(error)
    
    # 将错误信息重新放回队列
    for error in temp_list:
        error_queue.put(error)
