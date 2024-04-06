package test

import (
	"testing"

	"issc.io/mmitp/mmitp-core/utils"
)

func anyDuplicate(ids []string) bool {
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			return true
		}
		seen[id] = true
	}
	return false
}

func TestUuidGen(t *testing.T) {
	ids := make([]string, 0)
	for range [100000]int{} {
		id := utils.GenMqttCliId()
		ids = append(ids, id)
	}
	if anyDuplicate(ids) {
		t.Fatalf("duplicates found")
	}

}
