from sqlalchemy import Column, String, Float
from sqlalchemy.exc import SQLAlchemyError
from contextlib import contextmanager
from pythonp.apps.fortune_teller.models.transaction import db
import uuid


# 2. 定义模型
class TarotUser(db.Model):
    """
    Tarot用户表模型，包含 id 和 chat_id 以及 amount 积分字段。
    id: 主键，自增
    chat_id: 用户的 chat_id, 用于唯一标识用户
    amount: 积分，默认值 300
    """
    __tablename__ = 'tarot_users'
    id = Column(String(36),
                primary_key=True,
                default=lambda: str(uuid.uuid4()),
                unique=True,
                nullable=False)
    chat_id = Column(String(36), nullable=False)
    amount = Column(Float, nullable=False, default=300)

# 4. 事务上下文管理器（可选）
@contextmanager
def transaction_scope():
    """
    提供一个事务作用域，自动处理提交或回滚。
    """
    session = db.session()
    try:
        yield session
        session.commit()
    except Exception as e:
        session.rollback()
        raise e
    finally:
        session.close()


# 5. 更新积分函数
def update_amount(user_id, delta):
    """
    使用 SELECT ... FOR UPDATE 对指定 user_id 的用户进行行级锁定，
    并在同一事务中完成积分加减操作，保证数据原子性。
    """
    # 如果积分更新后可能为负，需要在逻辑内进行校验并抛出异常
    with transaction_scope() as session:
        try:
            # 1) 行级锁
            user = session.query(TarotUser)\
                          .filter(TarotUser.id == user_id)\
                          .with_for_update()\
                          .one()

            # 2) 业务校验：例如不能让积分小于 0
            if user.amount + delta < 0:
                raise ValueError("Insufficient amount")

            # 3) 更新积分
            user.amount += delta

            # 事务结束时自动提交
            return {
                "message": "amount updated successfully",
                "user_id": user.id,
                "new_amount": user.amount
            }
        except (SQLAlchemyError, ValueError) as e:
            # 出现任何异常都会在 transaction_scope 中触发回滚
            return { "error": str(e) }


# if __name__ == "__main__":
#     # 初始化数据库（若表不存在则创建）
#     init_db()

#     # 演示操作
#     with transaction_scope() as session:
#         # 添加一个用户, 初始积分100
#         new_user = TarotUser()
#         session.add(new_user)
#         # 此时会在 with 块结束时自动 commit
#     print(f"User created with ID: {new_user.id}, initial amount: {new_user.amount}")

#     # 1) 扣除20分
#     result = update_amount(new_user.id, -20)
#     print("更新积分 -20：", result)

#     # 2) 尝试扣除超出范围的积分 (比如 -200 分)
#     result = update_amount(new_user.id, -200)
#     print("尝试扣除 -200 分：", result)
