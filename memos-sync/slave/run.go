package slave

import (
	"fmt"
	"net/http"
	"taoey/memos-utils/memos-sync/dao"
	"taoey/memos-utils/memos-sync/util"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	// 获取memos 全量uid
	mux.HandleFunc("GET /memo/all_uids", func(w http.ResponseWriter, r *http.Request) {
		uids := GetAllMemoUids()
		w.Write([]byte(util.MustJsonStr(uids)))
	})
	// 获取memos uid 详情
	mux.HandleFunc("GET /memo/get_uid/{uid}", func(w http.ResponseWriter, r *http.Request) {
		uid := r.PathValue("uid")
		memoDetail := GetMemoDetailByUid(uid)
		w.Write([]byte(util.MustJsonStr(memoDetail)))
	})

	// 获取memos-关联情况
	mux.HandleFunc("GET /memo/relation", func(w http.ResponseWriter, r *http.Request) {
		result := GetMemoRelation()
		w.Write([]byte(util.MustJsonStr(result)))
	})

	// 获取memos-附件
	mux.HandleFunc("GET /memo/resource", func(w http.ResponseWriter, r *http.Request) {
		result := GetMemoResource()
		w.Write([]byte(util.MustJsonStr(result)))
	})

	ipPort := "0.0.0.0:8080"
	fmt.Println("slave running", ipPort)
	http.ListenAndServe(ipPort, mux)
}

func GetAllMemoUids() []string {
	db, _ := gorm.Open(sqlite.Open("/Users/th/Documents/memos/memos_prod.db"), &gorm.Config{})
	// 获取全量的memos UID 列表
	memoUids := []string{}
	sql := "select uid from memo"
	db.Raw(sql).Scan(&memoUids)
	return memoUids
}

func GetMemoDetailByUid(uid string) *dao.Memo {
	db, _ := gorm.Open(sqlite.Open("/Users/th/Documents/memos/memos_prod.db"), &gorm.Config{})

	memo := &dao.Memo{}
	db.Where("uid = ?", uid).First(&memo)
	return memo
}

func GetMemoRelation() []*dao.MemoRelationDTO {
	db, _ := gorm.Open(sqlite.Open("/Users/th/Documents/memos/memos_prod.db"), &gorm.Config{})
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
	db, _ := gorm.Open(sqlite.Open("/Users/th/Documents/memos/memos_prod.db"), &gorm.Config{})
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
