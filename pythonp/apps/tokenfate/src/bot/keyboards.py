from telegram import InlineKeyboardButton, InlineKeyboardMarkup
from typing import List, Dict, Any

class KeyboardFactory:
    """键盘布局工厂类"""
    
    def __init__(self, i18n_helper):
        self.i18n = i18n_helper
    
    def create_keyboard(self, layout_name: str, context: Any = None, **kwargs) -> InlineKeyboardMarkup:
        """
        创建键盘布局
        :param layout_name: 布局名称
        :param context: Telegram context
        :param kwargs: 其他参数
        :return: InlineKeyboardMarkup
        """
        if hasattr(self, f"_create_{layout_name}_keyboard"):
            keyboard_method = getattr(self, f"_create_{layout_name}_keyboard")
            return keyboard_method(context, **kwargs)
        raise ValueError(f"Unsupported keyboard layout: {layout_name}")
    
    def _create_start_keyboard(self, context: Any) -> InlineKeyboardMarkup:
        """创建开始界面的键盘布局"""
        keyboard = [
                [
                    InlineKeyboardButton(
                        self.i18n.get_button('launch', context=context),
                        callback_data='launch_to_reveal_button'
                    )
                ],
                [
                    InlineKeyboardButton(
                        self.i18n.get_button('community', context=context),
                        callback_data='join_community_button'
                    ),
                    # InlineKeyboardButton(
                    #     "🌐 Language/语言",
                    #     callback_data='switch_lang'
                    # )
                ],
                [
                    InlineKeyboardButton(
                        self.i18n.get_button('info', context=context),
                        callback_data='for_your_information_button'
                    )
                ],
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_info_keyboard(self, context: Any) -> InlineKeyboardMarkup:
        """创建开始界面的键盘布局"""
        keyboard = [
                [
                    InlineKeyboardButton(
                        self.i18n.get_button('launch', context=context),
                        callback_data='launch_to_reveal_button'
                    )
                ],
                [
                    InlineKeyboardButton(
                        self.i18n.get_button('community', context=context),
                        callback_data='join_community_button'
                    ),
                    # InlineKeyboardButton(
                    #     "🌐 Language/语言",
                    #     callback_data='switch_lang'
                    # )
                ],
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_risk_keyboard(self, context: Any) -> InlineKeyboardMarkup:
        """创建风险选择的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('role1', context=context),
                    callback_data='risk_role1'
                ),
                InlineKeyboardButton(
                    self.i18n.get_button('role2', context=context),
                    callback_data='risk_role2'
                )
            ],
            [
                InlineKeyboardButton(
                    self.i18n.get_button('role3', context=context),
                    callback_data='risk_role3'
                ),
                InlineKeyboardButton(
                    self.i18n.get_button('role4', context=context),
                    callback_data='risk_role4'
                )
            ]
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_unconnected_keyboard(self, context: Any=None, token: str=None) -> InlineKeyboardMarkup:
        """创建代币操作的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('connect', context=context),
                    callback_data='connect_wallet_button'
                )
            ],
            [
                InlineKeyboardButton(
                    self.i18n.get_button('reveal', context=context),
                    callback_data=f'reveal_fate:{token}'
                )
            ]
        ]
        return InlineKeyboardMarkup(keyboard)
    def _create_connected_keyboard(self, context: Any=None, token: str=None) -> InlineKeyboardMarkup:
        """创建代币操作的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('reveal', context=context),
                    callback_data=f'reveal_fate:{token}'
                )
            ]
        ]
        return InlineKeyboardMarkup(keyboard)