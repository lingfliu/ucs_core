package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/lingfliu/ucs_core/data/rtdb"
)

const HOST = "localhost:2366"

func main() {
	cli := rtdb.NewAgilorCli(HOST, "admin", "admin")

	sigRun, cancelRun := context.WithCancel(context.Background())

	go _task_rtdb(sigRun, cli)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	for {
		select {
		case <-s:
			cancelRun()
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func _task_rtdb(sigRun context.Context, cli *rtdb.AgilorCli) {
	cli.Open()
}
