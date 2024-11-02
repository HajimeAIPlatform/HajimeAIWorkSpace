import requests
from scrapy.crawler import CrawlerProcess
from hajime_scraper_lib.spiders.generic_spider import GenericSpider

def get_qian_function(response):
    print(f"Processing page: {response.url}")

    # 提取标题和 URL
    title = response.xpath('//title/text()').get()
    title = title.strip() if title else "无标题"
    url = response.url

    print(f"Title: {title}")
    print(f"URL: {url}")

    # 获取“抽签前的准备和注意事项”
    preparation_notice = response.xpath('//div[contains(text(), "抽签前的准备和注意事项")]/following-sibling::div[@class="lh25 f14 tleft gylqlist"][1]/descendant::text()').getall()
    preparation_content = format_content(preparation_notice, "抽签前的准备和注意事项")

    # 获取“灵签介绍”
    sign_intro = response.xpath('//div[contains(text(), "灵签介绍")]/following-sibling::div[@class="lh25 f14 tleft gylqlist"][1]/descendant::text()').getall()
    intro_content = format_content(sign_intro, "灵签介绍")

    # 获取“在线解签”总签数
    online_signs_section = response.xpath("//div[contains(@class, 'gylqtt') and contains(text(), '在线解签')]/following-sibling::div[1]/ul/li")
    online_signs_count = len(online_signs_section)

    print(f"在线解签总签数: {online_signs_count}")

    # 获取“签图和典故”中的图片 URL 和故事
    sign_image_url = response.xpath('//div[@class="qianzuo"]/img/@src').get()
    sign_story = response.xpath('//div[contains(@class, "gylqtt") and contains(text(), "签图和典故")]/following-sibling::div[@class="gylqlist"]/descendant::text()').getall()
    story_content = format_content(sign_story, "签图和典故")
    # 获取分类
    sign_category = url.split('/')[-2] if len(url.split('/')) > 4 else None  

    # 打印图片 URL 和故事
    if sign_image_url:
        print(f"签图 URL: {sign_image_url}")
    else:
        print("没有找到签图!")

    if url.endswith('.html'):

         # 获取“灵签概述”
        sign_overview = response.xpath('//div[contains(text(), "概述")]/following-sibling::div[@class="lq_m f14 tleft lh25"][1]/descendant::text()').getall()
        overview_content = format_content(sign_overview, "概述")

        # 查找包含 "解签" 两字的部分，并提取其后面的所有文本内容
        sign_explanation = response.xpath('//div[contains(text(), "解签")]/following-sibling::div//text()').getall()

        # 如果结构不同，考虑其他可能的 xpath
        if not sign_explanation:
            sign_explanation = response.xpath('//strong[contains(text(), "解签")]/following-sibling::div//text()').getall()

        # 合并提取到的文本内容
        sign_explanation = ''.join(sign_explanation).strip()
        # 获取“签图和典故”中的图片 URL 和故事
        sign_image_url = response.xpath('//div[@class="lq_m f14 tleft lh25"]/div/img/@src').get()

        # 构建解析后的数据
        main_parsed_data = {
            "sign_url": url,
            "sign_title": title,
            "sign_unscramble": sign_explanation,
            "sign_overview": overview_content,
            "sign_image_url": sign_image_url,
            "sign_category": sign_category,
        }
    else:
        # 构建解析后的数据
        main_parsed_data = {
            "sign_url": url,
            "sign_title": title,
            "sign_preparation_notice": preparation_content,
            "sign_intro": intro_content,
            "sign_total_signs": online_signs_count,
            "sign_image_url": sign_image_url,
            "sign_story": story_content,
            "sign_category": sign_category,
        }

    yield main_parsed_data
def format_content(content_list, label):
    """清洗并格式化内容"""
    if not content_list:
        print(f"No {label} found!")
        return f"无{label}"

    clean_content = [text.strip() for text in content_list if text.strip()]
    formatted_content = "\n\n".join(clean_content)
    print(f"{label}: {formatted_content}")
    return formatted_content

# 运行爬虫
def run_generic_spider(url, parse_function, domain_filter=None, crawl_all=False):
    process = CrawlerProcess(settings={
        'FEEDS': {
            'output.json': {
                'format': 'json',
                'encoding': 'utf8',
            },
        },
        'LOG_LEVEL': 'DEBUG',
        'DOWNLOADER_MIDDLEWARES': {
            'scrapy.downloadermiddlewares.httpproxy.HttpProxyMiddleware': 110,
            'hajime_scraper_lib.middlewares.ProxyMiddleware': 100,
        },
    })

    process.crawl(GenericSpider, 
                  start_url=url, 
                  parse_function=parse_function, 
                  domain_filter=domain_filter, 
                  crawl_all=crawl_all)
    process.start()

if __name__ == "__main__":
    url = 'https://www.zhouyi.cc/lingqian/'
    run_generic_spider(url, get_qian_function, domain_filter='https://www.zhouyi.cc/lingqian/', crawl_all=True)
