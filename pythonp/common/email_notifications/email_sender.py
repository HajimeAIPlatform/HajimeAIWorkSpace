# email_sender.py
import smtplib
import logging
from email.mime.text import MIMEText

def send_email(smtp_config, subject, body):

    msg = MIMEText(body)
    msg['Subject'] = subject
    msg['From'] = smtp_config['email_from']
    msg['To'] = smtp_config['email_to']

    try:
        with smtplib.SMTP_SSL(smtp_config['smtp_host'], smtp_config['smtp_port']) as server:
            server.login(smtp_config['smtp_user'], smtp_config['smtp_pass'])
            server.sendmail(smtp_config['email_from'], smtp_config['email_to'], msg.as_string())
        logging.info("邮件发送成功")
    except Exception as e:
        logging.info(f"发送邮件失败: {e}")