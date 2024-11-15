import json
import os
from pathlib import Path
from flask import session
from functools import wraps
from telegram import Update
from telegram.ext import CallbackContext

class I18nHelper:
    """国际化助手类"""
    
    def __init__(self, lang: str = 'zh'):
        self.messages = self._load_messages()
        self.lang = lang
    #     self.supported_languages = ['en', 'zh']
    
    def _load_messages(self):
        """加载语言配置文件"""
        config_path = 'pythonp/apps/tokenfate/static/assets/i18n_messages.json'
        default_path = 'static/assets/i18n_messages.json'

        # 检查文件是否存在
        if os.path.exists(config_path):
            path_to_use = config_path
        else:
            path_to_use = default_path
        try:
            with open(path_to_use, 'r', encoding='utf-8') as f:
                return json.load(f)
        except Exception as e:
            print(f"Error loading language file: {e}")
            return {}
    
    def reload_messages(self):
        """重新加载语言配置"""
        self.messages = self._load_messages()
    
    # def get_user_language(self, context: CallbackContext = None):
    #     """获取用户语言设置"""
    #     if context and 'language' in context.user_data:
    #         return context.user_data['language']
    #     return session.get('language', self.default_lang)
    
    # def set_user_language(self, lang: str, context: CallbackContext = None):
    #     """设置用户语言"""
    #     if lang not in self.supported_languages:
    #         return False
            
    #     if context:
    #         context.user_data['language'] = lang
    #     session['language'] = lang
    #     return True
    
    def get_dialog(self, key: str, lang: str = 'zh', context: CallbackContext = None):
        """获取对话文本"""
        return self.messages.get('dialogs', {}).get(self.lang, {}).get(key, '')
    
    def get_button(self, key: str, lang: str = 'zh', context: CallbackContext = None):
        """获取按钮文本"""
        return self.messages.get('buttons', {}).get(self.lang, {}).get(key, '')
