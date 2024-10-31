import json
import csv
import re


# 读取 JSON 文件
with open('output.json', 'r', encoding='utf-8') as file:
    data = json.load(file)

# 创建一个字典来存储每个 sign_category 的主项
category_dict = {}

# 定义汉字数字转换函数
def chinese_to_digit(chinese_num):
    chinese_numerals = {'一': 1, '二': 2, '三': 3, '四': 4, '五': 5, '六': 6, '七': 7, '八': 8, '九': 9, '十': 10}
    if chinese_num.isdigit():
        return int(chinese_num)
    
    if '十' in chinese_num:
        parts = chinese_num.split('十')
        if parts[0] == '':
            return 10 + chinese_numerals.get(parts[1], 0)
        elif parts[1] == '':
            return chinese_numerals[parts[0]] * 10
        else:
            return chinese_numerals[parts[0]] * 10 + chinese_numerals.get(parts[1], 0)
    else:
        return chinese_numerals.get(chinese_num, None)

# 定义提取签号的函数
def extract_sign_number(sign_title):
    match = re.search(r'第\s*([\d一二三四五六七八九十]+)\s*签', sign_title)
    if match:
        num_str = match.group(1)
        if num_str.isdigit():
            return int(num_str)
        else:
            return chinese_to_digit(num_str)
    return None

# 遍历数据并整合
for item in data:
    category = item['sign_category']
    item['sign_number'] = extract_sign_number(item['sign_title'])
    
    if item['sign_url'].endswith(f"/{category}/"):
        category_dict[category] = item
        category_dict[category]['related_signs'] = []

for item in data:
    category = item['sign_category']
    if not item['sign_url'].endswith(f"/{category}/"):
        if category in category_dict:
            category_dict[category]['related_signs'].append(item)

# 对 related_signs 进行排序
for category, item in category_dict.items():
    item['related_signs'] = sorted(item['related_signs'], key=lambda x: (x['sign_number'] if x['sign_number'] is not None else float('inf')))

merged_data = sorted(list(category_dict.values()), key=lambda x: (x['sign_number'] if x['sign_number'] is not None else float('inf')))

# 写入 JSON 文件
with open('merged_output_with_sign_numbers_sorted.json', 'w', encoding='utf-8') as file:
    json.dump(merged_data, file, ensure_ascii=False, indent=4)

# 导出为 CSV
with open('output.csv', 'w', newline='', encoding='utf-8') as csvfile:
    fieldnames = ['category', 'sign_url', 'sign_title', 'sign_number', 'sign_overview', 'sign_unscramble', 'sign_image_url']
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
    
    writer.writeheader()
    
    for item in merged_data:
        # 写主项
        writer.writerow({
            'category': item['sign_category'],
            'sign_url': item['sign_url'],
            'sign_title': item['sign_title'],
            'sign_number': item['sign_number'],
            'sign_overview': item.get('sign_overview', ''),
            'sign_unscramble': item.get('sign_unscramble', ''),
            'sign_image_url': item.get('sign_image_url', '')
        })
        
        # 写子项
        for related_sign in item['related_signs']:
            writer.writerow({
                'category': related_sign['sign_category'],
                'sign_url': related_sign['sign_url'],
                'sign_title': related_sign['sign_title'],
                'sign_number': related_sign['sign_number'],
                'sign_overview': related_sign.get('sign_overview', ''),
                'sign_unscramble': related_sign.get('sign_unscramble', ''),
                'sign_image_url': related_sign.get('sign_image_url', '')
            })

print("数据已导出为 CSV 文件！")
