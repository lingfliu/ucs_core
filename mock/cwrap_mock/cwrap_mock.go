package main

import (
	"github.com/lingfliu/ucs_core/ulog"
)

func main() {
	ulog.Config(ulog.LOG_LEVEL_INFO, "", false)
	// res := rtdb.Add(1, 2)
	// ulog.Log().I("main", fmt.Sprintf("res=%d", res))
	ulog.Log().I("main", "hello world")

}
