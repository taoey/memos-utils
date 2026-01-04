import requests
import json
import os
import base64
import mimetypes

from util.config import CONFIG_DICT

# --- 配置信息 ---
MEMOS_API_URL = CONFIG_DICT.get("memos_api_url")
ACCESS_TOKEN = "bearer " + CONFIG_DICT.get("access_token")

HEADERS = {"Authorization": ACCESS_TOKEN, "Content-Type": "application/json"}


# 更新笔记时间
# http://192.168.3.98:5000/update_create_time?memos_name=DK4emLtzPoexMKTwr2QCds&create_time=1763654400
def update_memos_create_time(name, create_time):
    url = f"{CONFIG_DICT.get('update_memos_create_time_server_domain')}/update_create_time?memos_name={name}&create_time={create_time}"
    resp = requests.get(url)
    return resp


# 上传附件
def upload_attachment(file_path):
    """
    使用 AttachmentService 上传文件
    API: POST /api/v1/attachments
    """
    url = f"{MEMOS_API_URL}/api/v1/attachments"

    if not os.path.exists(file_path):
        print(f"❌ 文件不存在: {file_path}")
        return None

    filename = os.path.basename(file_path)
    print(f"📤 正在上传附件: {filename} ...")

    try:
        mime_type, _ = mimetypes.guess_type(file_path)
        if not mime_type:
            mime_type = "application/octet-stream"  # 默认值
        # 以二进制方式读取文件
        with open(file_path, 'rb') as image_file:
            # 构造 multipart/form-data
            # Memos 后端通常识别 'file' 或 'content' 字段
            encoded_string = base64.b64encode(image_file.read()).decode('utf-8')
            data = {
                'name': filename,
                'filename': filename,
                "content": encoded_string,
                "type": mime_type,
            }
            # 发送请求
            response = requests.post(url, headers=HEADERS, json=data)

        response.raise_for_status()
        data = response.json()

        # 获取返回的资源标识符，通常在 'name' 字段中
        # 格式可能是 "attachments/123" 或 "resources/123"
        resource_name = data.get("name")
        print(f"✅ 上传成功! 资源名: {resource_name}")
        return resource_name

    except Exception as e:
        print(f"❌ 上传失败 [{filename}]: {e}")
        if 'response' in locals():
            print(f"服务器响应: {response.text}")
        return None


# 创建笔记
def create_memos(tags, content, creat_time=0, attachment_paths=[]):
    """创建一个新的备忘录

    Args:
        tags (list): 备忘录的标签列表
        content (str): 备忘录的内容
        creat_time (int, optional): 备忘录的创建时间戳. 如果为0的话默认使用当前时间，不为0的话使用指定时间
        attachment_paths (list, optional): 附件文件路径列表. Defaults to [].
    """
    data = {
        "content": content,
        "attachments": [],
    }
    if tags:
        tag_string = " ".join([f"#{tag}" for tag in tags])
        data["content"] = f"{tag_string}\n{content}"

    # 添加附件
    for file in attachment_paths:
        file_name = upload_attachment(file)
        if file_name is None:
            return
        data['attachments'].append({
            "name": file_name,
        })

    api_endpoint = f"{MEMOS_API_URL}/api/v1/memos"
    try:
        response = requests.post(api_endpoint, headers=HEADERS, data=json.dumps(data))

        if response.status_code == 200:
            memo_data = response.json()
            memo_name = memo_data.get("name")
            print("✅创建笔记成功", memo_data)
            if int(creat_time) > 0:
                update_memos_create_time(memo_name, str(creat_time))
        else:
            print(f"❌ 创建记录失败。状态码: {response.status_code}")
            print(f"错误信息: {response.text}")
    except requests.exceptions.RequestException as e:
        print(f"❌ 请求过程中发生错误: {e}")


# def main():
#     now = int(time.time())
#     create_memos(["python_test",'from_微博'],"这是通过API创建的测试记录",now,["aa.jpg"])
#     # create_memos(["python_test",'from_微博'],"这是通过API创建的测试记录")
#     return

# if __name__ == "__main__":
#     main()
