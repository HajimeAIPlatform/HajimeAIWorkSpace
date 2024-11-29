import uuid
import json
import logging
from typing import Dict
from os import getenv

from flask import Blueprint, jsonify, request, Response

from pythonp.apps.tokenfate.src.binance.utils import get_all_prices, process_recommendation
from pythonp.apps.tokenfate.dify_client import Client, models
from pythonp.apps.tokenfate.src.binance.schedule import get_random_dex_historical_prices

dify_api_key_workflow = getenv('DIFY_API_KEY_WORKFLOW')
dify_api_key_message = getenv('DIFY_API_KEY_MESSAGE')
dify_api_key_decode = getenv('DIFY_API_KEY_DECODE')
dify_api_key_2 = getenv('DIFY_API_KEy_2')
dify_api_base = getenv('DIFY_BASE_HOST')
if not dify_api_base or not dify_api_key_workflow or not dify_api_key_2 or not dify_api_key_message or not dify_api_key_decode:
    raise ValueError(
        "Dify API key and base host are not set in the environment")

# Initialize the client with your API key
client = Client(
    api_key=dify_api_key_workflow,
    api_base=dify_api_base,
)

client2 = Client(
    api_key=dify_api_key_2,
    api_base=dify_api_base,
)

client3 = Client(
    api_key=dify_api_key_message,
    api_base=dify_api_base,
)

client4 = Client(
    api_key=dify_api_key_decode,
    api_base=dify_api_base,
)

dify = Blueprint('dify', __name__)

def chat_blocking(data):
    try:
        user = str(uuid.uuid4())
        logging.info("Generated user ID: %s", user)

        today_market_data = get_random_dex_historical_prices()

        print(len(json.dumps(today_market_data)),'today_market_data')
        logging.info("Received data: %s", data)
        # Create a blocking chat request
        blocking_chat_req = models.ChatRequest(
            query=data.get("query"),
            inputs=data.get("inputs", {
                "today_market_data": json.dumps(today_market_data)
            }),
            user=str(data.get("user", user)),
            response_mode=models.ResponseMode.BLOCKING,
        )

        logging.info("Sending blocking chat request: %s", blocking_chat_req)

        # Send the chat message
        chat_response = client3.chat_messages(blocking_chat_req, timeout=60.)
        chat_response_dict = json.loads(
            json.dumps(chat_response,
                       default=lambda o: o.__dict__))  # Convert to dictionary

        # logging.info("Received chat response: %s", chat_response_dict)

        # Extract the answer from the chat response
        answer = chat_response_dict.get('answer', 'No answer found')
        logging.info("Answer: %s", answer)
        # 处理推荐结果
        return answer

    except Exception as e:
        logging.error("Error during chat_blocking: %s", e)
        return "An error occurred"

def chat_blocking_key_2(data):
    try:
        user = str(uuid.uuid4())
        logging.info("Generated user ID: %s", user)

        # Create a blocking chat request
        blocking_chat_req = models.ChatRequest(
            query=data.get("query"),
            inputs=data.get("inputs", {}),
            user=str(data.get("user", user)),
            response_mode=models.ResponseMode.BLOCKING,
        )

        logging.info("Sending blocking chat request: %s", blocking_chat_req)

        # Send the chat message
        chat_response = client2.chat_messages(blocking_chat_req, timeout=60.)
        chat_response_dict = json.loads(
            json.dumps(chat_response,
                       default=lambda o: o.__dict__))  # Convert to dictionary

        # logging.info("Received chat response: %s", chat_response_dict)

        # Extract the answer from the chat response
        answer = chat_response_dict.get('answer', 'No answer found')
        logging.info("Answer: %s", answer)
        # 处理推荐结果
        return answer

    except Exception as e:
        logging.error("Error during chat_blocking: %s", e)
        return "An error occurred"

def chat_workflow(data: Dict):
    try:
        # Generate a unique user ID if not provided
        user = str(uuid.uuid4())
        logging.info("User ID: %s", user)

        # Retrieve market data
        today_market_data = get_random_dex_historical_prices()
        serialized_market_data = json.dumps(today_market_data)

        # Prepare chat request with correct input format
        inputs: Dict = {
            "todayMarketData": serialized_market_data,
            "riskPreference": data.get("risk_preference", "neutral"),
        }

        # Prepare chat request
        blocking_chat_req = models.WorkflowsRunRequest(
            inputs=inputs,
            user=user,
            response_mode=models.ResponseMode.BLOCKING,
        )

        # Send the chat message and get response
        chat_response = client.run_workflows(blocking_chat_req, timeout=60)
        logging.info("Received chat response: %s", chat_response)
        # Check the status and handle errors
        status = chat_response.data.status
        error = chat_response.data.error

        if status == "succeeded" and error is None:
            # Extract the text from the response
            text = chat_response.data.outputs.get('text', 'No text found')
            logging.info("Text: %s", text)
            return text
        else:
            logging.error("Chat workflow failed with status: %s, error: %s", status, error)
            return "An error occurred during chat workflow"

    except Exception as e:
        logging.error("Error during chat_workflow: %s", e)
        return f"An unexpected error occurred,{e}"


def chat_decode(data):
    try:
        user = str(uuid.uuid4())
        logging.info("Generated user ID: %s", user)
        logging.info("Received data: %s", data)
        # Create a blocking chat request
        blocking_chat_req = models.ChatRequest(
            query=data.get("query", ""),
            inputs=data.get("inputs", {}),
            user=str(data.get("user", user)),
            response_mode=models.ResponseMode.BLOCKING,
        )

        logging.info("Sending blocking chat request: %s", blocking_chat_req)

        # Send the chat message
        chat_response = client4.chat_messages(blocking_chat_req, timeout=60)
        logging.info("Received chat response: %s", chat_response)
        chat_response_dict = json.loads(
            json.dumps(chat_response,
                       default=lambda o: o.__dict__))  # Convert to dictionary

        logging.info("Received chat response: %s", chat_response_dict)

        # Extract the answer from the chat response
        answer = chat_response_dict.get('answer', 'No answer found')
        logging.info("Answer: %s", answer)
        # 处理推荐结果
        return answer

    except Exception as e:
        logging.error("Error during chat_blocking: %s", e)
        return "An error occurred"

def chat_streaming(data):
    try:
        user = str(uuid.uuid4())
        logging.info("Generated user ID: %s", user)

        streaming_chat_req = models.ChatRequest(
            query=data.get("query"),
            inputs=data.get("inputs", {}),
            user=user,
            response_mode=models.ResponseMode.STREAMING,
        )

        logging.info("Sending streaming chat request: %s", streaming_chat_req)

        for chunk in client.chat_messages(streaming_chat_req, timeout=60.):
            chunk_dict = json.loads(
                json.dumps(
                    chunk,
                    default=lambda o: o.__dict__))  # Convert to dictionary
            answer = chunk_dict.get('answer', None)
            if answer is not None or answer != '':
                # logging.info("Received chunk of answer: %s", answer)
                yield answer

    except Exception as e:
        logging.error("Error during chat_streaming: %s", e)
        yield ""


@dify.route('/chat', methods=['POST'])
def handle_chat():
    data = request.json
    if not data or 'query' not in data:
        return jsonify({"error": "Invalid request, 'query' is required"}), 400

    response_mode = data.get("response_mode", models.ResponseMode.BLOCKING)

    if response_mode == models.ResponseMode.STREAMING:
        try:
            def generate():
                for chunk in chat_streaming(data):
                    yield f"data: {chunk}\n\n"
            return Response(generate(), content_type='text/event-stream')
        except Exception as e:
            logging.error("Error handling streaming chat request: %s", e)
            return jsonify({"error": "An error occurred while processing the request"}), 500
    else:
        try:
            chat_res = chat_blocking(data)
            result = process_recommendation(chat_res)
            if isinstance(result, tuple):
                token, action, amount = result
                return jsonify({"answer": chat_res, "recommendation": {
                    "token": token.upper(),
                    "action": action.lower(),
                    "amount": amount
                }}), 200

            return jsonify({"answer": chat_res}), 200
        except Exception as e:
            logging.error("Error handling blocking chat request: %s", e)
            return jsonify({"error": "An error occurred while processing the request"}), 500
