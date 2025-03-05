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
BOT_TOKEN="7578797243:AAG4gsM1m4Gn-fj9cZJ98yvbjyO4siSkzC4"
#BOT_TOKEN="7317644050:AAESKpy0z3-bMZrphdKRKJK925HaquaWQk4"

# BOT_TOKEN="7805651398:AAHCG40KeHJrqRqxXmlA_xDsEWwmMwOpnug"   # TokenFateFortuneCookieBot

# 使用 curl 命令设置 Telegram bot 的 webhook
curl -F "url=$WEBHOOK_URL" https://api.telegram.org/bot$BOT_TOKEN/setWebhook

# bash utils/set_webhook.sh  https://4a83-14-155-107-35.ngrok-free.app/telebot/bot/webhook
# Dev服务器:
# bash utils/set_webhook.sh https://hajimedev.pointer.ai/telebot/bot/webhook


# bash utils/set_webhook.sh https://d9ef-14-153-94-63.ngrok-free.app/telebot/bot_fortune_teller/webhook

# # company:
# export http_proxy="http://10.10.100.72:7897"
# export https_proxy="http://10.10.100.72:7897"

## home:
# export https_proxy="http://192.168.1.102:7897"
# export http_proxy="http://192.168.1.102:7897"

## 二级代理
# export http_proxy="http://192.168.1.100:3128"
# export https_proxy="http://192.168.1.100:3128"

## None
# export http_proxy=""
# export https_proxy=""

## 公网穿透
# cd ~/.config/ngrok
# ngrok start --config=ngrok.yml api