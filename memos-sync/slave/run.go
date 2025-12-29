package slave

import (
	"fmt"
	"net/http"
	"os"
	"taoey/memos-utils/memos-sync/dao"
	"taoey/memos-utils/memos-sync/util"

	"github.com/olebedev/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Config *config.Config //global config

var (
	ipPort         string
	accessPassword string
	dbFilePath     string
)

func InitConf() {
	pwd, _ := os.Getwd()
	configPath := pwd + "/config/slave.yml"
	Config, _ = config.ParseYamlFile(configPath)

	fmt.Println(util.MustJsonStr(Config))
	ipPort = Config.UString("ip_port")
	accessPassword = Config.UString("access_token")
	dbFilePath = Config.UString("db_filepath")
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Access-Token") != accessPassword {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Run() {
	// 加载配置文件
	InitConf()

	mux := http.NewServeMux()
	mux.Handle("GET /", auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})))
	// 获取memos 全量uid
	mux.Handle("GET /memo/all_uids", auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uids := GetAllMemoUids()
		w.Write([]byte(util.MustJsonStr(uids)))
	})))

	// 获取memos uid 详情
	mux.Handle("GET /memo/get_uid/{uid}", auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := r.PathValue("uid")
		memoDetail := GetMemoDetailByUid(uid)
		w.Write([]byte(util.MustJsonStr(memoDetail)))
	})))

	// 获取memos-关联情况
	mux.Handle("GET /memo/relation", auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := GetMemoRelation()
		w.Write([]byte(util.MustJsonStr(result)))
	})))

	// 获取memos-附件
	mux.HandleFunc("GET /memo/resource", func(w http.ResponseWriter, r *http.Request) {
		result := GetMemoResource()
		w.Write([]byte(util.MustJsonStr(result)))
	})
	fmt.Println("slave running", ipPort)
	err := http.ListenAndServe(ipPort, mux)
	fmt.Println(err)
}

func GetAllMemoUids() []string {
	db, _ := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	// 获取全量的memos UID 列表
	memoUids := []string{}
	sql := "select uid from memo"
	db.Raw(sql).Scan(&memoUids)
	return memoUids
}

func GetMemoDetailByUid(uid string) *dao.Memo {
	db, _ := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})

	memo := &dao.Memo{}
	db.Where("uid = ?", uid).First(&memo)
	return memo
}

func GetMemoRelation() []*dao.MemoRelationDTO {
	db, _ := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	sql := `
		select a.type,b.uid,c.uid related_memo_uid
		from memo_relation a  
		left join memo b on a.memo_id=b.id
		left join memo c on a.related_memo_id=c.id
	`
	result := []*dao.MemoRelationDTO{}
	db.Raw(sql).Scan(&result)
	return result
}

func GetMemoResource() []*dao.SlaveMemoResource {
	db, _ := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	resources := []*dao.Resource{}
	db.Find(&resources)

	slaveResources := []*dao.SlaveMemoResource{}
	for _, item := range resources {
		memo := &dao.Memo{}
		db.Where("id = ?", item.MemoID).First(&memo)
		slaveResources = append(slaveResources, &dao.SlaveMemoResource{
			Resource: *item,
			MemosUid: memo.UID,
		})
	}
	return slaveResources
}
