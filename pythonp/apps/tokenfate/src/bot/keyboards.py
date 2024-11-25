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
                        url="t.me/HajimeAI"
                    ),
                    InlineKeyboardButton(
                        self.i18n.get_button('info', context=context),
                        callback_data='for_your_information_button'
                    )
                ],
                [
                    InlineKeyboardButton(
                        self.i18n.get_button('aura', context=context),
                        callback_data='show_aura_rules'
                    )
                ],
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_info_keyboard(self, context: Any) -> InlineKeyboardMarkup:
        """创建详情界面的键盘布局"""
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
                        url="t.me/HajimeAI"
                    )
                ],
                [
                    InlineKeyboardButton(
                        self.i18n.get_button('aura', context=context),
                        callback_data='show_aura_rules'
                    )
                ],
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_risk_keyboard(self, context: Any=None, token: str=None, risk_target: str=None) -> InlineKeyboardMarkup:
        """创建风险选择的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('risk_role_Conservative', context=context),
                    callback_data=f'risk:role=Conservative&token={token}&risk_target={risk_target}',
                ),
                InlineKeyboardButton(
                    self.i18n.get_button('risk_role_Cautious', context=context),
                    callback_data=f'risk:role=Cautious&token={token}&risk_target={risk_target}',
                ),
                InlineKeyboardButton(
                    self.i18n.get_button('risk_role_Assertive', context=context),
                    callback_data=f'risk:role=Assertive&token={token}&risk_target={risk_target}',
                ),
                InlineKeyboardButton(
                    self.i18n.get_button('risk_role_Risk_Taker', context=context),
                    callback_data=f'risk:role=Risk_Taker&token={token}&risk_target={risk_target}',
                )
            ],
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_aura_keyboard(self, context: Any=None) -> InlineKeyboardMarkup:
        """创建积分规则的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('aura_action_wallet_connect', context=context),
                    callback_data='connect_wallet_button'
                ),
                InlineKeyboardButton(
                    self.i18n.get_button('aura_action_daily_checkin', context=context),
                    callback_data='aura_action_daily_checkin'
                ),
                InlineKeyboardButton(
                    self.i18n.get_button('aura_action_recommend_click', context=context),
                    callback_data='aura_action_recommend_click'
                ),
                InlineKeyboardButton(
                    self.i18n.get_button('aura_action_fate_reveal', context=context),
                    callback_data='launch_to_reveal_button'
                ),
            ],
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_unconnected_keyboard(self, context: Any=None, token: str=None) -> InlineKeyboardMarkup:
        """创建未连接钱包时的键盘布局"""
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
                    callback_data=f'reveal_fate:{token}:normal'
                )
            ]
        ]
        return InlineKeyboardMarkup(keyboard)
    def _create_connected_keyboard(self, context: Any=None, token: str=None) -> InlineKeyboardMarkup:
        """创建已连接钱包的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('reveal', context=context),
                    callback_data=f'reveal_fate:{token}:normal'
                )
            ]
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_launch_keyboard(self, context: Any=None) -> InlineKeyboardMarkup:
        """创建揭示的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('aura_action_fate_reveal', context=context),
                    callback_data='launch_to_reveal_button'
                )
            ]
        ]
        return InlineKeyboardMarkup(keyboard)
    
    def _create_lot_keyboard(self, context: Any=None, token: str=None, role: str=None) -> InlineKeyboardMarkup:
        """创建灵签的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('lot_decode', context=context),
                    callback_data=f'decode:{token}:{role}'
                )
            ]
        ]
        return InlineKeyboardMarkup(keyboard)

    def _create_lang_keyboard(self, context: Any=None) -> InlineKeyboardMarkup:
        """创建语言选择的键盘布局"""
        keyboard = [
            [
                InlineKeyboardButton(
                    self.i18n.get_button('lang_en', context=context),
                    callback_data='lang:en'
                )
            ],
            [
                InlineKeyboardButton(
                    self.i18n.get_button('lang_zh', context=context),
                    callback_data='lang:zh'
                )
            ],
            [
                InlineKeyboardButton(
                    self.i18n.get_button('lang_ja', context=context),
                    callback_data='lang:ja'
                )
            ]
        ]
        return InlineKeyboardMarkup(keyboard)