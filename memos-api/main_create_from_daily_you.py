import sqlite3
import re

from util.api import create_memos


def convert_hash_tags(text: str) -> str:
    """
    将文本中 【#...#】 格式的标签转换为 【#...】，
    其中 ... 必须是非空内容（至少一个字符）。
    特别地，【##】不会被转换。
    
    示例：
        "【#小记#】" → "【#小记】"
        "【##】"     → 保持不变
        "【#A1!你好#】" → "【#A1!你好】"
    """
    return re.sub(r'#(.+?)#', r'#\1 ', text)


mood_dict = {
    -2: "【😭 极度难过】",
    -1: "【😞 难过】",
    0: "【😐 一般】",
    1: "【😊 开心】",
    2: "【😁 非常开心】",
}


def query_data(db_path: str):
    """
    在 memo 表中查找 uid = target_uid 的记录，
    并将其 created_ts 和 updated_ts 更新为当前时间戳（秒）。
    
    :param db_path: SQLite 数据库文件路径
    :param target_uid: 要查找和更新的 uid
    :return: True if updated, False if not found
    """
    try:
        conn = sqlite3.connect(db_path)
        conn.row_factory = sqlite3.Row  # 允许通过列名访问
        cursor = conn.cursor()
        # 更新 created_ts 和 updated_ts
        sql = """
            select id,text,strftime('%s', time_create) AS time_create,mood from entries
        """
        cursor.execute(sql)
        rows = cursor.fetchall()
        for row in rows:
            id = row['id']
            text = row['text']
            text = convert_hash_tags(text)
            text = text.replace("## #", "###")
            time_create = row['time_create']
            mood = row['mood']

            # dt = datetime.strptime(time_create, "%Y-%m-%dT%H:%M:%S.%fZ")
            # timestamp = int(dt.timestamp())

            timestamp = time_create
            # print(id,timestamp,mood_dict.get(mood,""),text)

            # 通过ID查询图片
            sql_img = """
                select img_path from entry_images where entry_id = ?
            """
            cursor.execute(sql_img, (id, ))
            img_rows = cursor.fetchall()
            cur_image_files = []
            for img_row in img_rows:
                img_path = f"daily_backup/.images/{img_row['img_path']}"
                cur_image_files.append(img_path)

            # 创建笔记
            content = f"\n心情：{mood_dict.get(mood,'')}\n{text}"
            print(content)
            create_memos.create_memos(['from_daily-you-app'], content, timestamp, cur_image_files)
        return True

    except sqlite3.Error as e:
        print(f"SQLite error: {e}")
        return False
    finally:
        if conn:
            conn.close()


def main():
    db_path = "daily_backup/daily_you.db"
    query_data(db_path)


main()
