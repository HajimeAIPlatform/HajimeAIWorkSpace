import asyncio
import os
from threading import Thread, Event
import logging
import redis.asyncio as redis
from pythonp.common.email_notifications.email_sender import send_email 

class EmailMonitor:
    def __init__(self, smtp_config, check_interval=60*60):
        self.smtp_config = smtp_config
        self.check_interval = check_interval
        self.is_running = False
        self.monitor_thread = None
        self.stop_event = Event()
        self.redis_client = None  # Redis 客户端将在 _run_async_loop 中初始化

    async def _send_email_loop(self):
        """邮件发送循环"""
        while not self.stop_event.is_set():
            try:
                logging.info("邮件检查线程正在运行...")

                if not await self._is_exceptions_empty():
                    try:
                        subject, body = await self._create_email_content()
                        send_email(self.smtp_config, subject, body)
                        await self._clear_exceptions()
                        logging.info("邮件发送成功")
                    except Exception as e:
                        logging.error(f"邮件发送失败: {e}")

                logging.info(f"邮件检查线程开始休眠，等待 {self.check_interval} 秒后下一次检查...")
                await asyncio.sleep(self.check_interval)

            except asyncio.CancelledError:
                logging.info("邮件检查线程被取消")
                break
            except Exception as e:
                logging.error(f"邮件检查线程发生异常: {e}")

    def _run_async_loop(self):
        """在新线程中运行事件循环"""
        self.loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self.loop)
        try:
            # 初始化 Redis 客户端
            self.redis_client = redis.Redis(
                host=os.getenv("REDIS_HOST", "localhost"),
                port=int(os.getenv("REDIS_PORT", "6379")),
                decode_responses=True
            )
            self.loop.run_until_complete(self._send_email_loop())
        finally:
            if self.redis_client:
                self.redis_client.close()
                self.loop.run_until_complete(self.redis_client.wait_closed())
            self.loop.close()

    def start(self):
        """启动监控"""
        if self.is_running:
            logging.info("邮件监控已经在运行中")
            return
            
        self.is_running = True
        logging.info("邮件监控正在启动...")

        # 创建并启动新线程
        self.monitor_thread = Thread(target=self._run_async_loop, daemon=True)
        self.monitor_thread.start()
        logging.info("邮件监控已在新线程中启动")

    def stop(self):
        """停止监控"""
        if not self.is_running:
            logging.info("邮件监控未运行")
            return
            
        self.is_running = False
        self.stop_event.set()
        if self.monitor_thread and self.monitor_thread.is_alive():
            self.monitor_thread.join()
        logging.info("邮件监控已停止")

    async def _is_exceptions_empty(self):
        """检查异常列表是否为空"""
        llen = await self.redis_client.llen("exceptions")
        return llen == 0

    async def _create_email_content(self):
        """创建邮件内容"""
        exceptions = await self.redis_client.lrange("exceptions", 0, -1)
        subject = "TokenFate Error Report"
        body = "\n".join([exc for exc in exceptions])
        return subject, body

    async def _clear_exceptions(self):
        """清除所有存储的错误信息"""
        await self.redis_client.delete("exceptions")
        logging.info("Cleared all stored exceptions")

# 示例：如何使用 EmailMonitor 类
if __name__ == "__main__":
    smtp_config = {
        'server': 'smtp.example.com',
        'port': 587,
        'username': 'user@example.com',
        'password': 'password',
        'from_addr': 'user@example.com',
        'to_addrs': ['admin@example.com']
    }
    email_monitor = EmailMonitor(smtp_config, check_interval=60)
    email_monitor.start()

    try:
        # 模拟主程序运行一段时间
        import time
        time.sleep(300)
    finally:
        email_monitor.stop()