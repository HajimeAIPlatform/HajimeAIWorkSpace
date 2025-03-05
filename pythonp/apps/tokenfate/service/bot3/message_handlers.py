from pythonp.apps.tokenfate.service.dify.views import chat_tarot
from md2tgmd import escape

async def reply_chat_tarot(update, context):
    inputs = {
        "input": update.message.text,
    }
    answer = chat_tarot(inputs)
    await update.message.reply_photo(
        photo=answer["url"],
        caption=escape(answer["text"]),
        parse_mode="MarkdownV2",
    )