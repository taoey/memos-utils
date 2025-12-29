package master

import (
	"fmt"
	"log"
	"os"
	"taoey/memos-utils/memos-sync/dao"
	"taoey/memos-utils/memos-sync/util"

	"github.com/olebedev/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Config *config.Config //global config

var (
	accessPassword string
	dbFilePath     string
	url            string
)

func InitConf() {
	pwd, _ := os.Getwd()
	configPath := pwd + "/config/master.yml"
	Config, _ = config.ParseYamlFile(configPath)

	fmt.Println(util.MustJsonStr(Config))
	accessPassword = Config.UString("access_token")
	dbFilePath = Config.UString("db_filepath")
	url = Config.UString("url")
}

func Run() {
	// 初始化配置
	InitConf()

	c := util.NewHttpClient(url, accessPassword)
	// 1、同步memos
	var uids []string
	err := c.GetJSON("/memo/all_uids", nil, &uids)
	if err != nil {
		fmt.Println(err)
	}

	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印 SQL
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, uid := range uids {
		// 查询是否存在
		var count int64
		db.Model(&dao.Memo{}).Where("uid = ?", uid).Count(&count)
		if count > 0 {
			continue
		}
		// 入库
		memo := dao.Memo{}
		c.GetJSON("/memo/get_uid/"+uid, nil, &memo)
		memo.ID = 0
		// 插入数据库
		err := db.Create(&memo)
		if err != nil {
			fmt.Println(err)
		}
	}

	var memoRelations []*dao.MemoRelationDTO
	// 2、同步关联性，memos
	c.GetJSON("/memo/relation", nil, &memoRelations)
	for _, item := range memoRelations {
		relation := dao.MemoRelation{
			Type: item.Type,
		}
		// uid to id
		uid2idSql := `select id from memo where uid = ?`
		db.Raw(uid2idSql, item.UID).Scan(&relation.MemoID)
		db.Raw(uid2idSql, item.RelatedMemoUID).Scan(&relation.RelatedMemoID)

		var count int64
		db.Model(&dao.MemoRelation{}).
			Where("memo_id = ? and related_memo_id=? and type = ?",
				relation.MemoID, relation.RelatedMemoID, relation.Type).
			Count(&count)
		if count > 0 {
			continue
		}

		result := db.Create(&relation)
		fmt.Println(result.Error)
	}

	// 3、同步附件
	fmt.Println("开始同步附件......")
	var resources []*dao.SlaveMemoResource
	c.GetJSON("/memo/resource", nil, &resources)
	for _, item := range resources {
		// 判断是否已经存在
		if item.Resource.UID == "" {
			continue
		}
		var count int64
		db.Model(&dao.Resource{}).Where("uid = ?", item.Resource.UID).Count(&count)
		if count > 0 {
			continue
		}

		// 不存在则需要新增
		uid2idSql := `select id from memo where uid = ?`
		db.Raw(uid2idSql, item.MemosUid).Scan(&item.Resource.MemoID)

		item.Resource.ID = 0
		err := db.Create(&item.Resource)
		if err != nil {
			fmt.Println(err)
		}
	}
	// 关闭底层数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection closed")
}
