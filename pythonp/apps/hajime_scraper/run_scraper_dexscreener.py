import time
import csv
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
import undetected_chromedriver as uc


def get_dex_info(driver):
    print(f"Processing page: {driver.current_url}")

     # 检查复选框是否存在并点击
    try:
        # 等待复选框元素加载
        checkbox = WebDriverWait(driver, 20).until(
            EC.presence_of_element_located((By.CLASS_NAME, 'checkbox'))  # 查找复选框元素
        )
        
        # 检查复选框是否未被选中，如果是，则点击它
        if checkbox.is_selected():
            print("Checkbox is already selected.")
        else:
            # 如果复选框未选中，则模拟点击它
            ActionChains(driver).move_to_element(checkbox).click().perform()
            print("Checkbox clicked successfully")
            
    except Exception as e:
        print(f"Error finding or clicking the checkbox: {e}")

    # 等待页面加载数据
    time.sleep(5)
    

    # 获取表格标题
    headers = [
        "Rank","Token","Name","Icon","Price", "Age", "Txns", "Volume", "Makers", 
        "5M", "1H", "6H", "24H", "Liquidity", "MCAP"
    ]

    # 获取交易对信息
    pairs = driver.find_elements(By.XPATH, '//a[contains(@class, "ds-dex-table-row")]')
    print(f"Found {len(pairs)} pairs")

    # 打开 CSV 文件准备写入
    with open('dex_info.csv', mode='w', newline='', encoding='utf-8') as file:
        writer = csv.DictWriter(file, fieldnames=headers)
        writer.writeheader() 
        
        for pair in pairs:
            pair_data = {}

            # 提取每个交易对的各项数据
            try:
                rank = pair.find_element(By.XPATH, './/span[contains(@class, "ds-dex-table-row-badge-pair-no")]').text
            except:
                rank = "No Rank"
            try:
                icon_url = pair.find_element(By.XPATH, './/img[contains(@class, "ds-dex-table-row-chain-icon")]').get_attribute('src')
            except:
                icon_url = "No Icon"
            
            try:
                token_name = pair.find_element(By.XPATH, './/span[contains(@class, "ds-dex-table-row-base-token-name-text")]').text
            except:
                token_name = "No Token"

            try:
                pair_name = pair.find_element(By.XPATH, './/span[contains(@class, "ds-dex-table-row-base-token-symbol")]').text
                quote_token = pair.find_element(By.XPATH, './/span[contains(@class, "ds-dex-table-row-quote-token-symbol")]').text
                trading_pair = f"{pair_name}/{quote_token}" if pair_name and quote_token else "无交易对信息"
            except:
                trading_pair = "No Trading Pair"

            try:
                price = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-price")]').text
            except:
                price = "No Price"

            try:
                age = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-pair-age")]').text
            except:
                age = "No AGE"

            try:
                txns = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-txns")]').text
            except:
                txns = "No Txns"

            try:
                volume = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-volume")]').text
            except:
                volume = "No Volume"

            try:
                makers = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-makers")]').text
            except:
                makers = "No Makers"

            try:
                price_change_5m = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-price-change-m5")]').text
            except:
                price_change_5m = "No 5M"

            try:
                price_change_1h = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-price-change-h1")]').text
            except:
                price_change_1h = "No 1H"

            try:
                price_change_6h = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-price-change-h6")]').text
            except:
                price_change_6h = "No 6H"

            try:
                price_change_24h = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-price-change-h24")]').text
            except:
                price_change_24h = "No 24H"

            try:
                liquidity = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-liquidity")]').text
            except:
                liquidity = "No Liquidity"

            try:
                market_cap = pair.find_element(By.XPATH, './/div[contains(@class, "ds-dex-table-row-col-market-cap")]').text
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


# 启动 Selenium WebDriver
def run_selenium_spider(url):
    options = uc.ChromeOptions()
    options.add_argument('--headless')  # 无头模式
    options.add_argument('--no-sandbox')  # 避免沙盒问题
    options.add_argument('--disable-dev-shm-usage')  # 解决共享内存问题
    options.add_argument('--disable-gpu')
    options.add_argument("--verbose")
    options.binary_location = '/usr/bin/google-chrome'
    # 设置 ChromeDriver 路径
    service = Service('/usr/bin/chromedriver')  # 确保这是 ChromeDriver 的实际路径
    driver = uc.Chrome(service=service,options=options)

    try:
        print(f"Page title: {url}")

        driver.get(url)
        print(f"Page title: {url}{driver.title}")

        # 等待页面加载完成
        WebDriverWait(driver, 20).until(EC.title_is('DEX Screener'))
        print(f"Page title: {driver.title}")
        print(f"Page URL: {driver.current_url}")

        get_dex_info(driver)
    finally:
        driver.quit()

if __name__ == "__main__":
    url = 'https://dexscreener.com/?rankBy=trendingScoreH24&order=desc'  
    run_selenium_spider(url)
