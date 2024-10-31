import json
import re
from typing import List, Dict, Any

patterns = {

}

class DataProcess:
    def __init__(self):
        pass

    def get_related_signs(self, file_path = './signs.json') -> List[Dict[str, Any]]:
        with open(file_path, 'r', encoding='utf-8') as file:
            raws = json.load(file)
            all_signs = []
            for raw in raws:
                all_signs.extend(raw.get("related_signs", []))
        with open("./signs_thin.json", 'w', encoding='utf-8') as out_file:
            json.dump(all_signs, out_file, ensure_ascii=False, indent=4)
        return all_signs        

    def get_categories_and_write_samples(self, file_path = './signs.json') -> List[str]:
        with open(file_path, 'r', encoding='utf-8') as file:
            raws = json.load(file)
            categories = [ raw["sign_category"] for raw in raws]
            print(f'categories: {categories}, nums: {len(categories)}')
            samples = []
            for raw in raws:
                if "related_signs" in raw and raw["related_signs"]:
                    print(raw["related_signs"][0]["sign_category"])
                    samples.append(raw["related_signs"][0])
            # write samples data to samples.json
            with open('./samples.json', 'w', encoding='utf-8') as out_file:
                json.dump(samples, out_file, ensure_ascii=False, indent=4)
            return categories
            
    def write_to_json(self, file_path = "./signsProcessed.json") -> None:
        all_signs = self.load_json_data()
        processed_data = []
        for sign in all_signs:
            result = self.parse_text(sign)
            processed_data.append(result)
        try:
            with open(file_path, 'w', encoding='utf-8') as file:
                json.dump(processed_data, file, ensure_ascii=False, indent=4)
            print("Data successfully written to JSON file.")
        except IOError as e:
            print(f"Error writing to JSON file: {e}")
    def extract_available_data(self, file_path = "./signs_thin.json") -> List[Dict[str, Any]]:
        with open(file_path, 'r', encoding='utf-8') as file:
            signs = json.load(file)
            # 数据分类
            available_data = []
            unavailable_data = []
            for sign in signs:
                # 尝试解析
                result = self.parse_text(sign)
                if result:
                    available_data.append(result)
                else:
                    unavailable_data.append(sign)
        # 数据写入文件
        with open("./available_data.json", 'w', encoding='utf-8') as out_file:
            json.dump(available_data, out_file, ensure_ascii=False, indent=4)
        with open("./unavailable_data.json", 'w', encoding='utf-8') as out_file:
            json.dump(unavailable_data, out_file, ensure_ascii=False, indent=4)

    def parse_text(self, sign):
        title = sign["sign_title"]
        text = sign["sign_overview"]
        result = {}
        # Use regex patterns to extract the required sections

        # 解析签等级
        match_level = re.search(r'本签吉凶\n\n(.*?签)', text)
        if not match_level:
            return None
        else:
            result['sign_level'] = match_level.group(1)

        # 解析签诗
        match_level = re.search(r'签\s*诗\n\n(.*?)\n\n解\s*曰', text)
        if not match_level:
            return None
        else:
            result['sign_poem'] = match_level.group(1).strip()
        
        # 获取签名
        result['sign_name'] = title.strip().split()[0]

        # 解析解释
        # result['解曰'] = re.search(r'解\s*曰\n\n([\s\S]*?)\n\n圣\s*意', text).group(1).strip()
        # result['诗文解释'] = re.search(r'诗文解译\n\n([\s\S]*)$', text).group(1).strip()
        return result
    

if __name__ == '__main__':
    dp = DataProcess()
    dp.extract_available_data()
