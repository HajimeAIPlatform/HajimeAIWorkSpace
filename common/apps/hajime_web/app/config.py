"""
FastAPI server configuration
"""
from typing import List

# pylint: disable=too-few-public-methods

from decouple import config, Csv
from pydantic import BaseModel


class Settings(BaseModel):
    """Server config settings"""

    token_list: List[str] = config("ASSET_TOKENS", default=['USDT'], cast=Csv(str))
    invite_domain:str=config("INVITE_DOMAIN",default="",cast=str)
    openai_key:str=config("OPENAI_KEY",default="",cast=str)
    open_router_key:str=config("OPEN_ROUTER_KEY",default="",cast=str)
    use_open_router:str=config("USE_OPEN_ROUTER",default=False,cast=bool)
    db_save_key:str=config("DB_SAVE_KEY",default="",cast=str)
    private_key_path:str=config("PRIVATE_KEY_PATH",default="",cast=str)
    fs_store_dir:str=config("FS_STORE_PATH",default="",cast=str)
    file_domain:str=config("FILE_DOMAIN",default="",cast=str)

    system_agent_id:str=config("SYSTEM_AGENT_ID",default="",cast=str)
    app_key:str=config("APP_KEY",default="",cast=str)


CONFIG = Settings()
