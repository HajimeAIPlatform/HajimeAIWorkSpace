#!/bin/bash
###
 # @Description: 
 # @Author: Devin
 # @Date: 2024-12-19 17:11:10
### 

# 检查是否提供了 URL 参数
if [ -z "$1" ]; then
  echo "Usage: $0 <webhook_url>"
  exit 1
fi

# 设置 URL 参数
WEBHOOK_URL=$1
BOT_TOKEN="7886534816:AAF6GrZXaESSw0c7sFdWEbCrBQAd3_vJZ04" local tokenfate
# BOT_TOKEN="7317644050:AAESKpy0z3-bMZrphdKRKJK925HaquaWQk4" # prod tokenfate

# BOT_TOKEN="7805651398:AAHCG40KeHJrqRqxXmlA_xDsEWwmMwOpnug"   # TokenFateFortuneCookieBot

# 使用 curl 命令设置 Telegram bot 的 webhook
curl -F "url=$WEBHOOK_URL" https://api.telegram.org/bot$BOT_TOKEN/setWebhook

# bash utils/set_webhook.sh  https://8898-14-155-37-248.ngrok-free.app/telebot/bot/webhook
# Dev服务器:
# bash utils/set_webhook.sh https://hajimedev.pointer.ai/telebot/bot/webhook
# Prod服务器：
# bash utils/set_webhook.sh https://tokenfate.pointer.ai/telebot/bot/webhook

# bash utils/set_webhook.sh https://af65-119-136-89-174.ngrok-free.app/telebot/bot2/webhook