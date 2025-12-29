docker stop memos-master
rclone sync local-webdave:/assets /Users/th/Documents/memos-master/assets
go run main.go -process  master
docker restart memos-master