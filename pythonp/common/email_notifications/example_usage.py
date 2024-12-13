# main.py

import time
from error_collector import collect_error, start_error_checking

# 启动错误检查线程
start_error_checking()

# 模拟多个错误（包含重复）
try:
    raise ValueError("模拟错误1")
except Exception as e:
    collect_error(e)

time.sleep(10)  # 等待一段时间以模拟更多错误发生

try:
    raise ValueError("模拟错误1")  # 重复的错误
except Exception as e:
    collect_error(e)

try:
    raise ValueError("模拟错误2")
except Exception as e:
    collect_error(e)

# 模拟程序运行，保持主线程活跃以便观察效果
while True:
    time.sleep(1)
