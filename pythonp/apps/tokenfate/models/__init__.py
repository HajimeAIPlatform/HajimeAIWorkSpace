from flask import Flask
from string import Template
from models.transaction import db
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
    with app.app_context():
        create_all_tables()
