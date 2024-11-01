#!/bin/bash

# 检查是否提供了 URL 参数
if [ -z "$1" ]; then
  echo "Usage: $0 <webhook_url>"
  exit 1
fi

# 设置 URL 参数
WEBHOOK_URL=$1
#BOT_TOKEN="7437699352:AAEXiaMbmGwvxNaiZvd9gAwG5r4KuzSeMMA"
#BOT_TOKEN="7317644050:AAESKpy0z3-bMZrphdKRKJK925HaquaWQk4"

BOT_TOKEN="7476000600:AAGDllamIKdoMKu8fluDq15qVy4sCdiGgak"   # TokenFateFortuneCookieBot

# 使用 curl 命令设置 Telegram bot 的 webhook
curl -F "url=$WEBHOOK_URL" https://api.telegram.org/bot$BOT_TOKEN/setWebhook

# bash utils/set_webhook.sh https://121b-27-38-230-5.ngrok-free.app/telebot/bot/webhook
# Dev服务器:
# bash utils/set_webhook.sh https://hajimedev.pointer.ai/telebot/bot/webhook


# bash utils/set_webhook.sh https://af65-119-136-89-174.ngrok-free.app/telebot/bot2/webhook

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