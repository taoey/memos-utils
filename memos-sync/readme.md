
该工具用于同步不同memos实例的数据记录。

## 使用背景

【主memos】部署在家中局域网，【从memos】部署在公网云上。

有时可能访问不稳定，【从memos】进行临时记录，但是这些记录需要同步到【主memos】进行最终存档及维护。

## 主要功能
- 定时同步：主memos定时拉取从memos的相关数据记录进行合并
- 增量同步：从memos删除后，不影响已经同步到主memos的记录
- 同步机制：master节点pull slave 节点数据，因此master节点必须拥有从节点的访问权限


## 设计

从节点：
- 加载近n天的文章列表ID到内存中
- 通过ID拉取某个文章的全部内容（附件需要计算md5）
- 分享附件文件夹：安装webdav


主节点：
- 获取节点信息
- 安装rclone:同步附件 
brew install rclone

主要功能有三个表
- memo
- memo_relation
- reaction
- resource


其实memos的评论也是一条【memo记录】，是通过【memo_relation】表进行关联的




## 测试环境

创建master环境
```
docker run -d \
  --name memos-master \
  --restart unless-stopped \
  --cpus=0.7 \
  --memory=1024m \
  --memory-swap=1024m \
  -p 5231:5230 \
  -v /Users/th/Documents/memos-master:/var/opt/memos \
  docker.1ms.run/neosmemo/memos:0.25.2
```

rclone sync local-webdave:/assets /Users/th/Documents/memos-master/assets


主：
/var/opt/memos/assets/2025_12/2025-12-29/1766980737_Screenshot_20251228_235239.jpg

/Users/th/Documents/memos-master/assets/2025_12/2025-12-29/1766980737_Screenshot_20251228_235239.jpg