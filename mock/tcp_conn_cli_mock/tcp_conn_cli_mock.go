package main

import (
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/conn"
	"github.com/lingfliu/ucs_core/ulog"
)

func main() {

	ulog.Config(ulog.LOG_LEVEL_INFO, "", false)

	cfg := &conn.ConnCfg{
		RemoteAddr:     "127.0.0.1",
		Port:           12001,
		Class:          conn.CONN_CLASS_TCP,
		Timeout:        1000 * 1000 * 1000,
		TimeoutRw:      1000 * 1000 * 1000,
		KeepAlive:      true,
		ReconnectAfter: 1000 * 1000 * 1000,
	}

	cb := coder.NewCodebookFromJson("{}")
	// create a tcp client
	cli := conn.NewConnCli(cfg, cb)

	cli.HandleMsg = func(msg *coder.ZeroMsg) {
		ulog.Log().I("main", "received msg, class: "+strconv.Itoa(msg.Class))
	}

	cli.Start()
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	for {
		select {
		case <-s:
			cli.Close()
			return
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
