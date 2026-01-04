#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
@Time           : 2026-01-03 19:51
@Author         : tao
@Python Version : 3.13.3
@Desc           : None
"""
import sqlite3
import json
from util.api import create_memos
import re

WEIBO_PIC_DIR = "输入自己的微博图片目录"
WEIBO_DB_PATH = "输入自己的微博数据库路径"

conn = sqlite3.connect(WEIBO_DB_PATH)
cursor = conn.cursor()


sql = """
select id,tags,content,pic_url,create_time from post_item  order by id
"""

cursor.execute(sql)
result = cursor.fetchall()


def main():
    for row in result:
        id, tags_str, content, pic_url_str, create_time_str = row
        tags = json.loads(tags_str)
        pic_url = json.loads(pic_url_str)
        create_time = int(create_time_str)
        tags.append("from_weibo")

        for tag in tags:
            content = content.replace(f"#{tag}#", "")

        pic_file_path = []
        for pic in pic_url:
            file_path = WEIBO_PIC_DIR +"/"+ pic.split('/')[-1]
            pic_file_path.append(file_path)

        print("当前正在处理：", id)
        create_memos(tags, content, create_time, pic_file_path)

main()
