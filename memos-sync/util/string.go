package util

import (
	"encoding/json"
	"fmt"
)

func MustJsonStr(data any) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}
	return string(bytes)
}
