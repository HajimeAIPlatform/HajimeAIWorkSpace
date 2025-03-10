from PIL import Image
import os

class ImageToolKit:
    def __init__(self, project_root):
        self.project_root = project_root

    def get_image_path(self, image_name):
        # 构建图片的绝对路径
        image_path = os.path.join(self.project_root, 'static', 'images', image_name)
        return image_path

    def create_combined_image(self, image_names, output_name='combined_image.jpg'):
        image_paths = [self.get_image_path(name) for name in image_names]
        images = [Image.open(path) for path in image_paths]
        
        # 假设所有图片大小相同，获取单张图片的宽高
        width, height = images[0].size
        
        # 创建一个新的空白图像，宽度为4倍单张图片的宽度
        combined_image = Image.new('RGB', (width * len(images), height))
        
        # 将每张图片粘贴到新图像中
        for i, image in enumerate(images):
            combined_image.paste(image, (i * width, 0))
        
        # 保存合并后的图片
        combined_image_path = os.path.join(self.project_root, 'static', 'images', output_name)
        combined_image.save(combined_image_path)
        
        return combined_image_path

if __name__ == '__main__':
    # 使用示例
    project_root = 'pythonp/apps/fortune_teller'
    image_combiner = ImageToolKit(project_root)
    image_names = ["Preserver.png", "Stabilizer.png", "Enterpriser.png", "Adventurer.png"]
    combined_image_path = image_combiner.create_combined_image(image_names, output_name='risk_preference_combined.png')

    print(f'Combined image saved at: {combined_image_path}')
