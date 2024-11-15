import os

def get_assets_path(filename):
    config_path = 'pythonp/apps/tokenfate/static/assets'
    default_path = 'static/assets'

    # 检查文件是否存在
    if os.path.exists(config_path):
        path_to_use = config_path
    else:
        path_to_use = default_path

    return os.path.join(path_to_use, filename)

def get_images_path(filename):
    config_path = 'pythonp/apps/tokenfate/static/images/'
    default_path = 'static/images/'

    # 检查文件是否存在
    if os.path.exists(config_path):
        path_to_use = config_path
    else:
        path_to_use = default_path

    return os.path.join(path_to_use, filename)
