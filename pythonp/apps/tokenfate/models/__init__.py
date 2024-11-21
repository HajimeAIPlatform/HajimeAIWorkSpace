from flask import Flask
from string import Template
from pythonp.apps.tokenfate.models.transaction import db
import os


def get_connection_url():
    base_url = Template(
        "postgresql${engine_name}://" +
        f"{os.environ.get('POSTGRESQL_USER', 'telebot')}:" +
        f"{os.environ.get('POSTGRESQL_PASSWORD', 'telebot')}@" +
        f"{os.environ.get('POSTGRESQL_HOST', '10.10.101.126')}:" +
        f"{os.environ.get('POSTGRESQL_PORT', '5432')}/" +
        f"{os.environ.get('POSTGRESQL_DATABASE', 'telebot')}")
    return base_url.substitute(engine_name="+psycopg2")  # Use psycopg2 driver


def create_all_tables():
    """
    Create all tables in the models.
    """
    db.create_all()


def setup_db(app: Flask):
    """
    Set up the models for the Flask app.

    Parameters:
    - app (Flask): The Flask app object.

    Returns:
    None
    """
    if not app.config.get('SQLALCHEMY_DATABASE_URI', None):
        app.config['SQLALCHEMY_DATABASE_URI'] = get_connection_url()
    db.init_app(app)
    # 示例 (SQLAlchemy)
    # from sqlalchemy import create_engine
    # from sqlalchemy.orm import sessionmaker

    # engine = None
    # Session = None

    # def init_db():
    #     global engine, Session
    #     engine = create_engine('your_database_url', pool_size=10, max_overflow=20)
    #     Session = sessionmaker(bind=engine)

    # def get_session():
    #     if Session is None:
    #         init_db()
    #     return Session()
    
    # 创建数据库引擎并配置连接池
    # app.config['SQLALCHEMY_ENGINE_OPTIONS'] = {
    #     "pool_size": 10,  # 连接池大小
    #     "max_overflow": 5,  # 控制在连接池达到最大值后可以创建的连接数
    #     "pool_timeout": 30,  # 如果没有可用连接等待的时间
    #     "pool_recycle": 1800,  # 池中的连接回收时间，防止长时间未使用的连接失效
    #     "pool_pre_ping": True  # 使用预ping保持连接活跃
    # }
    
    # # 记录连接池配置
    # logging.info("Database connection pool configured successfully with options: %s", app.config['SQLALCHEMY_ENGINE_OPTIONS'])

    # # 检查连接池状态
    # with app.app_context():
    #     engine = db.get_engine()
    #     pool = engine.pool
    #     logging.info("Connection pool status: size=%d, overflow=%d, timeout=%d, recycle=%d, pre_ping=%s",
    #                 pool.size(), pool._overflow, pool._timeout, pool._recycle, pool._pre_ping)
    with app.app_context():
        db.create_all()