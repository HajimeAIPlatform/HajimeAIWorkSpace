import time
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy import Column, DateTime, String, BigInteger, Float, Boolean, ForeignKey,JSON, func, Integer
from datetime import datetime, timezone, timedelta
import uuid
import logging

db = SQLAlchemy()


def get_current_time(time_delta=8):
    """
    获取当前时间，存储数据库
    :param time_delta: 相对UTC的时间差，小时单位，可按地区查询
    :return: datetime
    """
    return datetime.now(tz=timezone(timedelta(hours=time_delta)))


class TonTransaction(db.Model):
    """
    TON wallet transaction table
    """
    __tablename__ = "ton_transaction"
    id = Column(String(36),
                primary_key=True,
                default=lambda: str(uuid.uuid4()),
                unique=True,
                nullable=False)
    chat_id = Column(BigInteger, nullable=False)
    user_id = Column(BigInteger, nullable=False)
    transaction_id = Column(String(36),
                            unique=True,
                            nullable=False,
                            default=lambda: str(uuid.uuid4()))
    side = Column(String(4), nullable=False)  # 'buy' or 'sell'
    symbol = Column(String(10), nullable=False)
    timestamp = Column(DateTime(timezone=True), default=get_current_time())
    status = Column(String(20), nullable=False)
    address = Column(String(100), nullable=False)
    amount = Column(Float, nullable=False)
    fee = Column(Float, nullable=False)
    transaction_pair_id = Column(String(36), nullable=False)
    trace_link = Column(String(255), nullable=False)  # 新增字段
    asset_id = Column(String(36), ForeignKey('user_asset.id'), nullable=False)
    paylink_id = Column(String(36), ForeignKey('paylink.id'), nullable=False)

    def to_dict(self):
        return {
            'id': str(self.id),
            'user_id': self.user_id,
            'chat_id': self.chat_id,
            'transaction_id': self.transaction_id,
            'side': self.side,
            'symbol': self.symbol,
            'fee': self.fee,
            'timestamp': self.timestamp.isoformat() + 'Z',
            'status': self.status,
            'address': self.address,
            'amount': self.amount,
            'transaction_pair_id': self.transaction_pair_id,
            'trace_link': self.trace_link,  # 返回 trace_link
            'asset_id': self.asset_id,
            'paylink_id': self.paylink_id
        }


class BinanceTransaction(db.Model):
    """
    Binance transaction table
    """
    __tablename__ = "binance_transaction"
    id = Column(String(36),
                primary_key=True,
                default=lambda: str(uuid.uuid4()),
                unique=True,
                nullable=False)
    order_id = Column(String(36), unique=True, nullable=False)
    type = Column(String(36), nullable=False)  # 'market' or 'limit'
    side = Column(String(36), nullable=False)  # 'buy' or 'sell'
    timestamp = Column(BigInteger, nullable=False)
    status = Column(String(36), nullable=False)
    symbol = Column(String(36), nullable=False)
    amount = Column(Float, nullable=False)
    cummulative_quote_qty = Column(Float, nullable=False)
    fee = Column(Float, nullable=False)
    transaction_pair_id = Column(String(36), nullable=False)
    full_data = Column(JSON, nullable=False)

    def to_dict(self):
        return {
            'id': str(self.id),
            'order_id': self.order_id,
            'type': self.type,
            'side': self.side,
            'timestamp': self.timestamp,
            'status': self.status,
            'symbol': self.symbol,
            'amount': self.amount,
            'cummulative_quote_qty': self.cummulative_quote_qty,
            'fee': self.fee,
            'transaction_pair_id': self.transaction_pair_id,
        }


class UserAsset(db.Model):
    """
    User asset table
    """
    __tablename__ = "user_asset"
    id = Column(String(36),
                primary_key=True,
                default=lambda: str(uuid.uuid4()),
                unique=True,
                nullable=False)
    user_id = Column(BigInteger, nullable=False)
    symbol = Column(String(10), nullable=False)  # 币种
    amount = Column(Float, nullable=False)

    def to_dict(self):
        return {
            'id': str(self.id),
            'user_id': self.user_id,
            'symbol': self.symbol,
            'amount': self.amount
        }

    @classmethod
    def get_amount_by_user_id_and_symbol(cls, user_id, symbol):
        """
        Query the amount based on user_id and symbol.

        :param user_id: The user ID to query.
        :param symbol: The asset symbol to query.
        :return: The amount if found, else None.
        """
        symbol = symbol.upper()
        asset = db.session.query(cls).filter_by(user_id=user_id,
                                                symbol=symbol).first()
        if asset:
            return asset.amount
        return None

    @classmethod
    def get_assets_by_user_id(cls, user_id):
        """
        Get all assets for a given user_id.

        :param user_id: The user ID to query.
        :return: A dictionary with asset symbols as keys and their amounts as values.
        """
        assets = db.session.query(cls).filter_by(user_id=user_id).all()
        asset_dict = {asset.symbol: asset.amount for asset in assets}
        return asset_dict

class Paylink(db.Model):
    """
    Paylink table
    """
    __tablename__ = "paylink"
    id = Column(String(36),
                primary_key=True,
                default=lambda: str(uuid.uuid4()),
                unique=True,
                nullable=False)
    chat_id = Column(String(36), nullable=False)
    amount = Column(Float, nullable=False)
    trace_link = Column(String(255), nullable=False)

    def to_dict(self):
        return {
            'id': str(self.id),
            'chat_id': self.chat_id,
            'amount': self.amount,
            'trace_link': self.trace_link
        }


def save_paylink_to_db(chat_id, amount, trace_link):
    paylink = Paylink(chat_id=chat_id, amount=amount, trace_link=trace_link)
    db.session.add(paylink)
    db.session.commit()
    return paylink.id


def save_ton_transaction_to_db(user_id, chat_id, symbol, side, status, address,
                               amount, trace_link, fee):
    symbol = symbol.upper()
    paylink_id = save_paylink_to_db(chat_id, amount, trace_link)
    asset_id = save_user_asset_to_db(user_id, symbol, amount, side)
    transaction_pair_id = str(uuid.uuid4())
    ton_transaction = TonTransaction(
        user_id=user_id,
        chat_id=chat_id,
        symbol=symbol,
        side=side,
        status=status,
        address=address,
        amount=amount,
        transaction_pair_id=transaction_pair_id,
        trace_link=trace_link,  # 保存 trace_link
        paylink_id=paylink_id,
        asset_id=asset_id,
        fee=fee)
    db.session.add(ton_transaction)
    db.session.commit()

    return transaction_pair_id


def save_binance_transaction_to_db(side, status, symbol, amount,cummulative_quote_qty, fee, type,
                                   transaction_pair_id, order_id, timestamp,
                                   full_data):
    symbol = symbol.upper()
    binance_transaction = BinanceTransaction(
        type=type,
        status=status,
        symbol=symbol,
        amount=amount,
        cummulative_quote_qty=cummulative_quote_qty,
        transaction_pair_id=transaction_pair_id,
        fee=fee,
        side=side,
        order_id=order_id,
        timestamp=timestamp,
        full_data=full_data)
    db.session.add(binance_transaction)
    db.session.commit()


def save_user_asset_to_db(user_id, symbol, amount, side):
    symbol = symbol.upper()
    user_asset = UserAsset.query.filter_by(user_id=user_id,
                                           symbol=symbol).first()
    if user_asset:
        if side == 'BUY':
            user_asset.amount += amount
        elif side == 'SELL':
            user_asset.amount -= amount
    else:
        user_asset = UserAsset(user_id=user_id, symbol=symbol, amount=amount)
    db.session.add(user_asset)
    db.session.commit()
    return user_asset.id


class PointsLog(db.Model):
    __tablename__ = 'points_log'

    id = Column(Integer, primary_key=True)
    user_id = Column(BigInteger, nullable=False)
    current_points = Column(Integer, nullable=False)
    change_amount = Column(Integer, nullable=False)
    balance_after_change = Column(Integer, nullable=False)
    description = Column(String, nullable=True)
    timestamp = Column(DateTime, default=func.now(), nullable=False)


class InsufficientPointsError(Exception):
    """Exception raised for errors in the point deduction process."""
    
    def __init__(self, user_id, attempted_points, current_points, message="Insufficient points."):
        self.user_id = user_id
        self.attempted_points = attempted_points
        self.current_points = current_points
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return (f"User {self.user_id} has insufficient points. "
                f"Attempted to deduct {self.attempted_points}, "
                f"but only {self.current_points} available.")


class UserPoints(db.Model):
    """
    User points table
    """
    __tablename__ = "user_points"
    id = Column(String(36), primary_key=True, default=lambda: str(uuid.uuid4()), unique=True, nullable=False)
    user_id = Column(BigInteger, nullable=False, unique=True)
    points = Column(Integer, nullable=False, default=0)
    language = Column(String(2), nullable=False, default='') 
    daily_recommended_points = Column(Integer, nullable=False, default=0)
    last_reset_date = Column(DateTime, nullable=False, default=datetime.today().date())

    def to_dict(self):
        return {
            'id': str(self.id),
            'user_id': self.user_id,
            'points': self.points,
            'language': self.language,
            'daily_recommended_points': self.daily_recommended_points,
            'last_reset_date': self.last_reset_date
        }

    @classmethod
    def get_points_by_user_id(cls, user_id):
        user_points = db.session.query(cls).filter_by(user_id=user_id).first()
        if user_points:
            return user_points.points
        return 0

    @classmethod
    def update_points_by_user_id(cls, user_id, points, description=""):
        try:
            user_points = db.session.query(cls).with_for_update().filter_by(user_id=user_id).first()
            if user_points:
                current_points = user_points.points
                new_points = current_points + points
                if new_points < 0:
                    raise InsufficientPointsError(user_id, points, current_points)
                user_points.points = new_points
            else:
                if points < 0:
                    raise InsufficientPointsError(user_id, points, 0)
                user_points = UserPoints(user_id=user_id, points=points)
                db.session.add(user_points)
                current_points = 0

            # Log the transaction
            points_log = PointsLog(
                user_id=user_id,
                current_points=current_points,
                change_amount=points,
                balance_after_change=user_points.points,
                description=description
            )
            db.session.add(points_log)

            db.session.commit()
            return True
        except (SQLAlchemyError, InsufficientPointsError) as e:
            logging.error(f"Error updating user points: {e}")
            db.session.rollback()
            return False
    
    @classmethod
    def get_language_by_user_id(cls, user_id):
        user_points = db.session.query(cls).filter_by(user_id=user_id).first()
        if user_points:
            return user_points.language
        return ''
    
    @classmethod
    def update_language_by_user_id(cls, user_id, language):
        try:
            user_points = db.session.query(cls).with_for_update().filter_by(user_id=user_id).first()
            if user_points:
                user_points.language = language
            else:
                user_points = UserPoints(user_id=user_id, language=language)
                db.session.add(user_points)

            db.session.commit()
            return True
        except SQLAlchemyError as e:
            logging.error(f"Error updating user language: {e}")
            db.session.rollback()
            return False
        
    @classmethod
    def check_daily_recommended_points(cls, user_id):
        try:
            user_points = db.session.query(cls).with_for_update().filter_by(user_id=user_id).first()
            if user_points:
                today = datetime.today().date()
                if user_points.last_reset_date != today:
                    user_points.daily_recommended_points = 0
                    user_points.last_reset_date = today
                    db.session.commit()
                
                # 检查点击推荐积分是否超过上限
                if user_points.daily_recommended_points < 50:
                    return True
                else:
                    return False
            else:
                return True  # 新用户可以直接增加积分
        except SQLAlchemyError as e:
            logging.error(f"Error checking daily recommended points: {e}")
            db.session.rollback()
            return False
    @classmethod
    def update_daily_recommended_points(cls, user_id, points=10):
        try:
            user_points = db.session.query(cls).with_for_update().filter_by(user_id=user_id).first()
            if user_points:
                user_points.daily_recommended_points += points
                db.session.commit()
                return True
            else:
            # 新用户首次点击推荐
                user_points = UserPoints(user_id=user_id, daily_recommended_points=points, last_reset_date=datetime.today().date())
                db.session.add(user_points)
                db.session.commit()
                return True
        except SQLAlchemyError as e:
            logging.error(f"Error updating daily recommended points: {e}")
            db.session.rollback()
            return False


#     transaction_pair_id = save_ton_transaction_to_db(
#         user_id=user_id,
#         chat_id=chat_id,
#         trade_type=trade_type,
#         status=status,
#         symbol=symbol,
#         address=address,
#         trace_link=trace_link,
#         asset_id=None,
#         transaction_id=str(uuid.uuid4()),
#         amount=amount
#     )
#     save_binance_transaction_to_db(
#         trade_type=trade_type,
#         status=status,
#         symbol=symbol,
#         amount=amount,
#         transaction_pair_id=transaction_pair_id,
#         trace_link=trace_link
#     )
#     save_user_asset_to_db(user_id, symbol, amount if trade_type == 'buy' else -amount)
