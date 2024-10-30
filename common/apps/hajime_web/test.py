import requests
# 定义API的认证信息
import pymysql

# 打开数据库连接
db = pymysql.connect(host='localhost', user='root', passwd='samos1688', port=3306,
                     database="hajime")
print('连接成功！')
cursor = db.cursor()


# SQL 查询语句
sql = "SELECT * FROM wallet"
cursor.execute(sql)
# 获取所有记录列表
results = cursor.fetchall()
print(results)
for row in results:
    print(row)

    # 打印结果
    print('数据查询成功！')
cursor = db.cursor()

sql = "SELECT * FROM `order`"
cursor.execute(sql)
# 获取所有记录列表
results = cursor.fetchall()
for row in results:
    print(row)
    # 打印结果
    print('数据查询成功！')

db.close()

