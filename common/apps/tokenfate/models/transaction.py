import time
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import Column, DateTime, String, BigInteger, Float, Boolean, ForeignKey,JSON
from datetime import datetime, timezone, timedelta
import uuid

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
