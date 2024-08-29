package test

import (
	"testing"

	"github.com/lingfliu/ucs_core/utils"
)

func TestRandInt(t *testing.T) {
	t.Log("TestRandInt")
	v := utils.RandInt64(0, 100)
	if v < 0 || v > 100 {
		t.Fatalf("RandInt64 failed")
	}
}
