import time
import csv
from zenrows import ZenRowsClient
from bs4 import BeautifulSoup
import os

# 初始化ZenRowsClient
client = ZenRowsClient("31692dad13a2f39f41427fe581f9cfb01d604833")
url = "https://dexscreener.com/?rankBy=trendingScoreH24&order=desc" 
params = {"js_render":"true"}

# 使用ZenRowsClient获取网页内容
response = client.get(url, params=params)

# 检查响应状态码
if response.status_code == 200:
    
    # 获取表格标题
    headers = [
        "Rank", "Token", "Name", "Icon", "Price", "Age", "Txns", "Volume", "Makers", 
        "5M", "1H", "6H", "24H", "Liquidity", "MCAP"
    ]

    # 获取交易对信息
    soup = BeautifulSoup(response.text, 'html.parser')
    pairs = soup.find_all('a', class_="ds-dex-table-row")
    print(f"Found {len(pairs)} pairs")

    file_path = "/tmp/dex_info.csv"
    # 打开 CSV 文件准备写入
    with open(file_path, mode='w', newline='', encoding='utf-8') as file:
        writer = csv.DictWriter(file, fieldnames=headers)
        writer.writeheader() 
        
        for pair in pairs:
            pair_data = {}

            # 使用BeautifulSoup提取每个交易对的各项数据
            try:
                rank = pair.find('span', class_="ds-dex-table-row-badge-pair-no").text
            except:
                rank = "No Rank"
            try:
                icon_url = pair.find('img', class_="ds-dex-table-row-chain-icon")['src']
            except:
                icon_url = "No Icon"
            
            try:
                token_name = pair.find('span', class_="ds-dex-table-row-base-token-name-text").text
            except:
                token_name = "No Token"

            try:
                pair_name = pair.find('span', class_="ds-dex-table-row-base-token-symbol").text
                quote_token = pair.find('span', class_="ds-dex-table-row-quote-token-symbol").text
                trading_pair = f"{pair_name}/{quote_token}" if pair_name and quote_token else "无交易对信息"
            except:
                trading_pair = "No Trading Pair"

            try:
                price = pair.find('div', class_="ds-dex-table-row-col-price").text
            except:
                price = "No Price"

            try:
                age = pair.find('div', class_="ds-dex-table-row-col-pair-age").text
            except:
                age = "No AGE"

            try:
                txns = pair.find('div', class_="ds-dex-table-row-col-txns").text
            except:
                txns = "No Txns"

            try:
                volume = pair.find('div', class_="ds-dex-table-row-col-volume").text
            except:
                volume = "No Volume"

            try:
                makers = pair.find('div', class_="ds-dex-table-row-col-makers").text
            except:
                makers = "No Makers"

            try:
                price_change_5m = pair.find('div', class_="ds-dex-table-row-col-price-change-m5").text
            except:
                price_change_5m = "No 5M"

            try:
                price_change_1h = pair.find('div', class_="ds-dex-table-row-col-price-change-h1").text
            except:
                price_change_1h = "No 1H"

            try:
                price_change_6h = pair.find('div', class_="ds-dex-table-row-col-price-change-h6").text
            except:
                price_change_6h = "No 6H"

            try:
                price_change_24h = pair.find('div', class_="ds-dex-table-row-col-price-change-h24").text
            except:
                price_change_24h = "No 24H"

            try:
                liquidity = pair.find('div', class_="ds-dex-table-row-col-liquidity").text
            except:
                liquidity = "No Liquidity"

            try:
                market_cap = pair.find('div', class_="ds-dex-table-row-col-market-cap").text
            except:
                market_cap = "No Market Cap"

            # 填充数据
            pair_data = {
                "Rank": rank,
                "Icon": icon_url,
                "Name": token_name,
                "Token": trading_pair,
                "Price": price,
                "Age": age,
                "Txns": txns,
                "Volume": volume,
                "Makers": makers,
                "5M": price_change_5m,
                "1H": price_change_1h,
                "6H": price_change_6h,
                "24H": price_change_24h,
                "Liquidity": liquidity,
                "MCAP": market_cap
            }

            # 写入数据到 CSV
            writer.writerow(pair_data)
        print(f"Data has been written to {file_path}")
else:
    print(f"Failed to retrieve page, status code: {response.status_code}")
