import scrapy

class GenericSpider(scrapy.Spider):
    name = "generic_spider"

    def __init__(self, start_url, parse_function, crawl_all=False, domain_filter=None, *args, **kwargs):
        super(GenericSpider, self).__init__(*args, **kwargs)
        self.start_urls = [start_url]
        self.parse_function = parse_function
        self.crawl_all = crawl_all
        self.domain_filter = domain_filter  # 用于过滤链接，确保只爬取指定域名或路径下的页面

    def parse(self, response):
        """解析页面并递归爬取"""
        # 调用传入的解析函数来处理页面内容
        yield from self.parse_function(response)

        # 如果启用了 crawl_all 模式，则递归爬取页面中的所有链接
        if self.crawl_all:
            links = response.xpath('//a/@href').getall()
            for link in links:
                # 如果链接是相对链接，将其转化为绝对链接
                if link.startswith('/'):
                    next_page = response.urljoin(link)
                else:
                    next_page = link

                # 过滤掉不符合条件的链接（确保在指定域名或路径下）
                if self.domain_filter and not next_page.startswith(self.domain_filter):
                    continue  # 跳过不符合条件的链接

                # 发起对子页面的请求
                yield scrapy.Request(next_page, callback=self.parse)
