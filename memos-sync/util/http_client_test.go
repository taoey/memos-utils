package util_test

import (
	"fmt"
	"taoey/memos-utils/memos-sync/util"
	"testing"
)

func TestClient_GetJSON(t *testing.T) {
	c := util.NewHttpClient("http://localhost:8080")
	var resp any
	err := c.GetJSON("/memo/all_uids", nil, &resp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
