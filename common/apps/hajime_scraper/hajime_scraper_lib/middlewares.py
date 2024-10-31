import requests
import random
import base64

class ProxyMiddleware:
    def __init__(self):
        self.proxy_api_url = "http://api1.ydaili.cn/tools/BUnlimitedApi.ashx?key=6E65E4973ADB34005D73B4989B7BF7CEAD2F70EC0041766C&action=BUnlimited&qty=100&orderNum=SH20241016190539760&isp=&format=json"
        self.proxies = []  
        self.fetch_proxies()  
    def fetch_proxies(self):
        """从代理池 API 获取代理列表"""
        try:
            response = requests.get(self.proxy_api_url)
            if response.status_code == 200:
                data = response.json()
                print(data)
                if data.get('status') == 'success':
                    # 解析并存储代理
                    self.proxies = [proxy['IP'] for proxy in data['data']]
                    print(f"成功获取代理: {self.proxies}")
                else:
                    print(f"获取代理失败: {data.get('message', 'Unknown error')}")
            else:
                print(f"代理 API 请求失败，状态码: {response.status_code}")
        except Exception as e:
            print(f"获取代理时发生异常: {e}")

    def process_request(self, request, spider):
        """在每个请求之前为请求设置代理"""

        if not self.proxies:
            self.fetch_proxies()  # 如果没有代理，重新获取

        if self.proxies:
            proxy = random.choice(self.proxies)  # 随机选择一个代理
            request.meta['proxy'] = f"http://{proxy}"
            print(f"使用代理: {proxy}")
        else:
            print("没有可用代理")

    def process_exception(self, request, exception, spider):
        if 'proxy' in request.meta:
            proxy = request.meta['proxy']
            print(f"代理 {proxy} 连接失败，移除并尝试重新获取代理...")
            if proxy in self.proxies:
                self.proxies.remove(proxy)  # 移除无效的代理
            self.fetch_proxies()  # 获取新的代理
