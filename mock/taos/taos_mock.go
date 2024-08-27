package main

import (
	"fmt"
	"time"

	"github.com/lingfliu/ucs_core/data/rtdb"
)

func main() {
	cli := &rtdb.TaosCli{
		Host:     "localhost:6030",
		Username: "root",
		Password: "taosdata",
	}

	go _task_connect_test(cli)
	for {
		time.Sleep(1 * time.Second)
	}
}

func _task_connect_test(cli *rtdb.TaosCli) {
	cli.Open()
	cli.CreateTable("ucs", "eval_demo", "ts timestamp, val float")
	for i := 0; i < 100; i++ {
		cli.Insert("ucs", "eval_demo", "ts, val", fmt.Sprintf("now, %d", i))
	}
	defer cli.Close()
}
