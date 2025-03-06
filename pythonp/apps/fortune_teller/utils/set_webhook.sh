#!/bin/bash
###
 # @Description: 
 # @Author: Devin
 # @Date: 2025-03-05 15:32:22
### 
# 检查是否提供了 URL 参数
if [ -z "$1" ]; then
  echo "Usage: $0 <webhook_url>"
  exit 1
fi

# 设置 URL 参数
WEBHOOK_URL=$1
BOT_TOKEN="7578797243:AAG4gsM1m4Gn-fj9cZJ98yvbjyO4siSkzC4"

# 使用 curl 命令设置 Telegram bot 的 webhook
curl -F "url=$WEBHOOK_URL" https://api.telegram.org/bot$BOT_TOKEN/setWebhook


# fortune teller
# BOT_TOKEN="7578797243:AAG4gsM1m4Gn-fj9cZJ98yvbjyO4siSkzC4"
# bash utils/set_webhook.sh https://d9ef-14-153-94-63.ngrok-free.app/fortune_teller/bot_fortune_teller/webhook
# bash utils/set_webhook.sh https://tokenfate.pointer.ai/fortune_teller/bot_fortune_teller/webhook
