from flask import Flask, request, jsonify, render_template, abort
import sqlite3
import time

app = Flask(__name__)

#输入自己的db路径
DB_PATH = "/root/memos/memos/memos_prod.db"

def update_memo_timestamps(db_path: str, target_uid: str, create_time):
    """
    在 memo 表中查找 uid = target_uid 的记录，
    并将其 created_ts 和 updated_ts 更新为当前时间戳（秒）。
    
    :param db_path: SQLite 数据库文件路径
    :param target_uid: 要查找和更新的 uid
    :return: True if updated, False if not found
    """
    # 获取当前 Unix 时间戳（整数，单位：秒）
    current_ts = int(create_time)
    now = int(time.time())

    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()

        # 检查是否存在该 uid
        cursor.execute("SELECT 1 FROM memo WHERE uid = ?", (target_uid, ))
        if cursor.fetchone() is None:
            print(f"UID '{target_uid}' not found in memo table.")
            return False

        # 更新 created_ts 和 updated_ts
        cursor.execute(
            """
            UPDATE memo 
            SET created_ts = ?, updated_ts = ? 
            WHERE uid = ?
        """, (current_ts, now, target_uid))

        conn.commit()
        print(f"Successfully updated timestamps for UID '{target_uid}'.")
        return True

    except sqlite3.Error as e:
        print(f"SQLite error: {e}")
        return False
    finally:
        if conn:
            conn.close()


@app.route('/update_create_time', methods=['GET'])
def update_create_time():
    create_time = request.args.get('create_time', '')
    memos_name = request.args.get('memos_name', '')
    memos_name = memos_name.replace("memos/", "")

    print("参数", create_time, memos_name)

    if not create_time or not memos_name:
        abort(400, description="Missing 'create_time' or 'memos_name' parameter")

    update_memo_timestamps(DB_PATH, memos_name, create_time)

    resp = {"msg": "ok"}
    return jsonify(resp)


if __name__ == '__main__':
    # 运行 Flask 应用 在0.0.0.0
    app.run(host='0.0.0.0', port=6000, debug=True)
