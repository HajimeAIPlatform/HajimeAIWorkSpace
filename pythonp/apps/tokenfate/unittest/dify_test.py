from dotenv import load_dotenv

load_dotenv()

import uuid
import json
import logging
from os import getenv
from typing import Dict

from pydantic import ValidationError

from pythonp.apps.tokenfate.src.binance.schedule import get_random_usdt_historical_prices
from pythonp.apps.tokenfate.dify_client import Client, models

dify_api_key = getenv('DIFY_API_KEY')
dify_api_base = getenv('DIFY_BASE_HOST')
if not dify_api_base or not dify_api_key:
    raise ValueError(
        "Dify API key and base host are not set in the environment")

# Initialize the client with your API key
client = Client(
    api_key=dify_api_key,
    api_base=dify_api_base,
)


def chat_workflow(data: Dict):
    try:
        # Generate a unique user ID if not provided
        user = data.get("user", str(uuid.uuid4()))
        logging.info("User ID: %s", user)

        # Retrieve market data
        today_market_data = get_random_usdt_historical_prices()
        serialized_market_data = json.dumps(today_market_data)

        # Prepare chat request with correct input format
        inputs: Dict = {
            "todayMarketData": serialized_market_data,
        }

        # Prepare chat request
        blocking_chat_req = models.WorkflowsRunRequest(
            inputs=data.get("inputs", inputs),
            user=user,
            response_mode=models.ResponseMode.BLOCKING,
        )

        # Send the chat message and get response
        chat_response = client.run_workflows(blocking_chat_req, timeout=60)

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

if __name__ == "__main__":
    print(chat_workflow({}))
