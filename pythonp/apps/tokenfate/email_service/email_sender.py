import time
import logging
import asyncio
from threading import Thread
from datetime import datetime
from pythonp.apps.tokenfate.src.ton.tc_storage import ExceptionStorage
from pythonp.common.email_notifications.email_sender import send_email 

class EmailMonitor:
    def __init__(self, smtp_config, check_minutes=6):
        """
        初始化邮件监控器
        smtp_config: 邮件配置信息
        check_minutes: 检查间隔（分钟）
        """
        self.smtp_config = smtp_config
        # self.check_interval = check_minutes * 60  # 转换为秒
        self.check_interval = 1  # test用时
        self.storage = ExceptionStorage()
        self.is_running = False
        self.monitor_thread = None
        self.loop = None

    async def _create_email_content(self):
        """准备邮件内容"""
        exceptions = await self.storage.get_exceptions()
        body = "\n".join(exc.get('message', '') for exc in exceptions)
        subject = f"错误汇总 - {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}"
        return subject, body

    async def _send_email_loop(self):
        """邮件发送循环"""
        while self.is_running:
            try:
                print("邮件检查线程正在运行...")
                
                # 如果有错误信息，就发送邮件
                if not await self.storage.is_empty():
                    try:
                        subject, body = await self._create_email_content()
                        send_email(self.smtp_config, subject, body)
                        await self.storage.clear_exceptions()
                        print("邮件发送成功")
                    except Exception as e:
                        print(f"邮件发送失败: {e}")
                
                # 休眠指定时间
                print("邮件检查线程开始休眠，等待下一次检查...")
                await asyncio.sleep(self.check_interval)
            except asyncio.CancelledError:
                print("邮件检查线程被取消")
                break
            except Exception as e:
                print(f"邮件检查线程发生异常: {e}")

    def _run_async_loop(self):
        """在新线程中运行事件循环"""
        self.loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self.loop)
        try:
            self.loop.run_until_complete(self._send_email_loop())
        finally:
            self.loop.close()
            
    def start(self):
        """启动监控"""
        if self.is_running:
            print("邮件监控已经在运行中")
            return
            
        self.is_running = True
        print("邮件监控正在启动...")
        
        # 创建并启动新线程
        self.monitor_thread = Thread(target=self._run_async_loop, daemon=True)
        self.monitor_thread.start()
        print("邮件监控已在新线程中启动")
            
    def stop(self):
        """停止监控"""
        if not self.is_running:
            print("邮件监控已经停止")
            return
            
        print("正在停止邮件监控...")
        self.is_running = False
        
        # 等待线程结束
        if self.monitor_thread and self.monitor_thread.is_alive():
            self.monitor_thread.join(timeout=5)
            print("邮件监控线程已结束")
        
        # 关闭事件循环
        if self.loop and self.loop.is_running():
            self.loop.stop()
            print("事件循环已停止")

# 使用示例
if __name__ == "__main__":
    # 邮件配置
    smtp_config = {
        "host": "smtp.example.com",
        "port": 587,
        "username": "your_email@example.com",
        "password": "your_password",
        "sender": "sender@example.com",
        "recipients": ["recipient@example.com"]
    }
    
    # 创建并启动监控器
    monitor = EmailMonitor(smtp_config)
    monitor.start()
    
    # 保持主程序运行一段时间
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        monitor.stop()