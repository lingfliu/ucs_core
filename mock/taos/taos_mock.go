package main

import (
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
	cli.CreateSTable("ucs", "eval_demo", "ts timestamp, val float", "node_id string, offset int")
	for i := 0; i < 100; i++ {
		cli.Insert("eval_demo", []string{"ts", "val", "node_id", "offset"}, []string{"now", "3.14", "node1"})
	}
	defer cli.Close()
}
