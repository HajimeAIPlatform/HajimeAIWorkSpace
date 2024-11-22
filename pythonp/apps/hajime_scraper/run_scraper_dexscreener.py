import json
from bs4 import BeautifulSoup
import requests

DEFAULT_API_KEY = "0d0b52c02a37a758af91d9d83856f094"
DEFAULT_URL = "https://dexscreener.com/?rankBy=trendingScoreH24&order=desc"
DEFAULT_OUTPUT_PATH = "/tmp/dex_info.json"

def fetch_dex_data(api_key = DEFAULT_API_KEY, url = DEFAULT_URL, output_path="/tmp/dex_info.json"):
    """
    参数：
        api_key (str): ScraperAPI的API密钥。
        url (str): 需要抓取数据的目标URL。
        output_path (str): 保存JSON文件的路径，默认为/tmp/dex_info.json。
    
    返回：
        str: 保存的JSON文件路径（成功）或错误信息（失败）。
    """
    # ScraperAPI请求的参数
    payload = {'api_key': api_key, 'url': url}
    response = requests.get('https://api.scraperapi.com/', params=payload)

    # 检查响应状态码
    if response.status_code != 200:
        return f"Failed to retrieve page, status code: {response.status_code}"

    # 获取交易对信息
    soup = BeautifulSoup(response.text, 'html.parser')
    pairs = soup.find_all('a', class_="ds-dex-table-row")
    print(f"Found {len(pairs)} pairs")

    # 准备JSON数据
    data_list = []

    for pair in pairs:
        pair_data = {}

        # 使用BeautifulSoup提取每个交易对的各项数据
        try:
            rank = pair.find('span', class_="ds-dex-table-row-badge-pair-no").text
        except:
            rank = "No Rank"
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
            volume = pair.find('div', class_="ds-dex-table-row-col-volume").text
        except:
            volume = "No Volume"
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
            "Token": trading_pair,
            "Price": price,
            "Volume": volume,
            "5M": price_change_5m,
            "1H": price_change_1h,
            "6H": price_change_6h,
            "24H": price_change_24h,
            "Liquidity": liquidity,
            "MCAP": market_cap
        }

        data_list.append(pair_data)

    # 保存为JSON文件
    with open(output_path, mode='w', encoding='utf-8') as json_file:
        json.dump(data_list, json_file, ensure_ascii=False, indent=4)
    print(f"Data has been written to {output_path}")
    return output_path


if __name__ == "__main__":
    # 定义参数
    api_key = "0d0b52c02a37a758af91d9d83856f094"
    url = "https://dexscreener.com/?rankBy=trendingScoreH24&order=desc"
    output_path = "/tmp/dex_info.json"

    # 调用函数
    result = fetch_dex_data(api_key, url, output_path)
    
    # 打印结果
    if result.endswith(".json"):
        print(f"Data has been written to {result}")
    else:
        print(result)