import unittest
from pythonp.apps.tokenfate.service.ton.tc_storage import ExceptionStorage

class TestExceptionStorage(unittest.TestCase):
    def setUp(self):
        self.storage = ExceptionStorage()
    
    async def test_store_exception(self):
        # 测试将异常信息存储到Redis
        exception_message = "Test exception"
        await self.storage.store_exception(exception_message)

        # 验证异常信息是否存储成功
        exceptions = await self.storage.get_exceptions()
        self.assertEqual(len(exceptions), 1)
        self.assertEqual(exceptions[0]["message"], exception_message)

    async def test_get_exceptions(self):
        # 测试获取所有异常信息
        await self.storage.store_exception("Exception 1")
        await self.storage.store_exception("Exception 2")

        exceptions = await self.storage.get_exceptions()
        self.assertEqual(len(exceptions), 2)
        self.assertEqual(exceptions[0]["message"], "Exception 1")
        self.assertEqual(exceptions[1]["message"], "Exception 2")

    async def test_clear_exceptions(self):
        # 测试清除所有异常信息
        await self.storage.store_exception("Error to be cleared")

        await self.storage.clear_exceptions()

        exceptions = await self.storage.get_exceptions()
        self.assertEqual(len(exceptions), 0)
