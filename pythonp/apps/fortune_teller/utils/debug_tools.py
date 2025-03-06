import sys
import traceback

def get_user_friendly_error_info(exc_info=None):
    """从异常信息中提取用户代码的第一个栈帧并格式化为简洁的字符串。
    
    参数:
        exc_info (tuple): 包含异常类型、值和追踪对象的元组，默认为当前异常信息。
        
    返回:
        str: 格式化的异常信息字符串。
    """
    if exc_info is None:
        exc_info = sys.exc_info()
    
    exc_type, exc_value, exc_tb = exc_info
    tb = traceback.extract_tb(exc_tb)
    
    # 找到最后一个非标准库/第三方库的栈帧
    user_frames = [frame for frame in tb 
                   if "site-packages" not in frame.filename and "lib" not in frame.filename]
    
    last_frame = user_frames[-1] if user_frames else None
    
    if last_frame:
        path_parts = last_frame.filename.split('/')
        short_path = '/'.join(path_parts[-6:])
        error_message = f"{short_path}:{last_frame.name}:{last_frame.lineno} - {str(exc_value)}"
    else:
        error_message = str(exc_value)
    return error_message