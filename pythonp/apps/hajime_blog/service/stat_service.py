import pymysql

from pythonp.apps.hajime_blog.db.models import UserDeposit

# 打开数据库连接
db = pymysql.connect(host='localhost', user='root', passwd='samos1688', port=3306,
                     database="hajime")
print('连接成功！')
cursor = db.cursor()

class StatService:


    @classmethod
    async def stat(cls):
        items = await UserDeposit.find({})
        for item in items:
            address = item.sender
            cursor = db.cursor()
            sql = "SELECT * FROM user_deposit_list WHERE address = '%s'" % address
