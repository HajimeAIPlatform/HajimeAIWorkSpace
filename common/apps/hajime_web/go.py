import asyncio
import time
from datetime import datetime

import requests

from app.db.schemas import SolScanModel

api_key = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOjE3MTI3NzIwNzk2NjUsImVtYWlsIjoibm92YUBoYWppbWUuYWkiLCJhY3Rpb24iOiJ0b2tlbi1hcGkiLCJpYXQiOjE3MTI3NzIwNzl9.yIGBrunwuNwooWDnPen2nf9U-JKdnhdi-nY5JAwdT5E"
async def get_transfer(account, offset=0, limit=200, fromTime=0):
    headers = {
        "token": api_key
    }
    # url = f"https://pro-api.solscan.io/v1.0/account/splTransfers?account={account}&limit={limit}&offset={offset}"
    # if fromTime > 0:
    #     now = int(time.time())
    #     url = f"{url}"

    url = f"https://pro-api.solscan.io/v1.0/account/splTransfers?account={account}&limit={limit}&offset={offset}"
    if fromTime > 0:
        now = int(time.time())
        url = f"{url}&fromTime={fromTime}&toTime={now}"

        dt_object_from_ms = datetime.fromtimestamp(fromTime)

        now2 = datetime.now()
        formatted_now = now2.strftime("%Y-%m-%d %H:%M:%S")
        formatted_from = dt_object_from_ms.strftime("%Y-%m-%d %H:%M:%S")
        print(formatted_from, formatted_now)

    print(url)
    response = requests.get(url, headers=headers)

    # print(response.text)
    # print(response.text)
    try:
        data = response.json()
        # print(data)
        scan = SolScanModel(**data)
        # print("Total:", scan.total)
        for item in scan.data:
            print(item.blockTime)
            dt_object_from_ms = datetime.fromtimestamp(item.blockTime)
            formatted_from = dt_object_from_ms.strftime("%Y-%m-%d %H:%M:%S")
            print(formatted_from)

        return len(scan.data)
    except:
        print("json_decode error", response.text)
        return 0


async def check_transfer():


        account_list = [
            "BfWYgztHDqrvnf1RXGofDR49JPi7BHbkCmdGDyqSHtKe",
            # "Gg8avohYTZ9G4skbqXJbtqzc99NmkTJfAhSomELWQSsh",
            # "9npi4xTUNBwWPKaCtcWDuVVR9zsrqmWRQqLzNkeQLujG",
            # "ALvmTrNzPuyJKzsHGZZBPhjCF39cXXJGhDzBwSguCGXF"
        ]

        last_check_at = 1715773744
        for account in account_list:
            limit = 50
            offset = 0
            while True:
                num = await get_transfer(account, offset, limit, last_check_at)

                break
                if num > 0:
                    offset += limit
                    print("offset,limit:", offset, limit)
                else:
                    print("done")
                    break

if __name__ == "__main__":
    try:
        asyncio.run(check_transfer())
    except KeyboardInterrupt:
        print("async io")
        pass
