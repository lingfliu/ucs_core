package utils

import "encoding/json"

func CheckJson(jstr string) bool {
	var v any
	err := json.Unmarshal([]byte(jstr), v)
	if err != nil {
		return false
	} else {
		return true
	}
}
