# email_sender.py

import time
import smtplib
from email.mime.text import MIMEText
from threading import Thread
from datetime import datetime
from pythonp.common.email_notifications.email_global import error_queue, send_interval_minutes

def send_email(smtp_config, subject, body):

    msg = MIMEText(body)
    msg['Subject'] = subject
    msg['From'] = smtp_config['email_from']
    msg['To'] = smtp_config['email_to']

    try:
        with smtplib.SMTP_SSL(smtp_config['smtp_host'], smtp_config['smtp_port']) as server:
            server.login(smtp_config['smtp_user'], smtp_config['smtp_pass'])
            server.sendmail(smtp_config['email_from'], smtp_config['email_to'], msg.as_string())
        print("邮件发送成功")
    except Exception as e:
        print(f"发送邮件失败: {e}")

def print_queue_contents():
    # 打印当前队列中的所有错误信息
    print("Current errors in queue:")
    temp_list = []
    while not error_queue.empty():
        error = error_queue.get()
        temp_list.append(error)  # 将错误信息存储到临时列表中，以便后续使用
        print(error)
    
    # 将错误信息重新放回队列
    for error in temp_list:
        error_queue.put(error)


def prepare_email_body():
    body_lines = []
    while not error_queue.empty():  # 从队列中获取所有错误信息
        error = error_queue.get()
        body_lines.append(
            f"{error[0]} - {error[1]} "
            f"(File: {error[2]}, Function: {error[3]}, Line: {error[4]})"
        )
    return "\n".join(body_lines)

def email_sender_thread(smtp_config):
    while True:
        print("邮件发送线程正在休眠...")
        time.sleep(send_interval_minutes)  # 等待指定时间间隔
        print(f"Email sender thread accessing queue ID: {id(error_queue)}")
        if not error_queue.empty():  # 如果队列中有错误信息需要发送
            last_send_time = datetime.now()
            body = prepare_email_body()
            subject = f"错误汇总 - {last_send_time.strftime('%Y-%m-%d %H:%M:%S')}"
            
            try:
                send_email(smtp_config, subject, body)  # 调用发送邮件函数
            except Exception as e:
                print(f"邮件发送失败: {e}")


def start_email_sender(smtp_config):
    # 启动邮件发送线程并传递SMTP配置
    Thread(target=email_sender_thread, args=(smtp_config,), daemon=True).start()
