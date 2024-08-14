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

	cliConnCfg := &conn.ConnCfg{
		RemoteAddr: "127.0.0.1",
		Port:       12001,
		Class:      conn.CONN_CLASS_TCP,

		Timeout:        1000 * 1000 * 1000,
		TimeoutRw:      1000 * 1000 * 1000,
		KeepAlive:      true,
		ReconnectAfter: 1000 * 1000 * 1000,
	}

	cb := coder.NewCodebookFromJson("{}")
	srv := conn.NewConnSrv(cliConnCfg, cb)
	srv.MsgTimeout = 1000 * 1000 * 1000

	srv.MsgHandler = func(cc *conn.ConnCli, msg *coder.ZeroMsg) {
		ulog.Log().I("main", "msg from cli: "+cc.Conn.GetRemoteAddr()+" msg class: "+strconv.Itoa(msg.Class))
	}

	go srv.Start()
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	for {
		select {
		case <-s:
			srv.Stop()
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
