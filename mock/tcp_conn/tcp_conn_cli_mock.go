package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/conn"
	"github.com/lingfliu/ucs_core/ulog"
)

func main() {

	ulog.Config(ulog.LOG_LEVEL_INFO, "./log.log", false)

	cfg := &conn.ConnCfg{
		RemoteAddr:     "127.0.0.1",
		Port:           12001,
		Class:          conn.CONN_CLASS_TCP,
		Timeout:        1000,
		TimeoutRw:      1000,
		KeepAlive:      true,
		ReconnectAfter: 1000,
	}

	// create a tcp client
	cli := conn.NewCli(cfg, "")

	//expose callbacks
	cli.OnRecvBytes = func(b []byte) {
		ulog.Log().I("tcp_conn_cli_mock", "recv bytes at "+time.Now().String()+" "+string(b))
	}

	cli.OnReq = func(m coder.UMsg) *coder.UMsg {
		ulog.Log().I("tcp_conn_cli_mock", "req msg received at "+time.Now().String())
		return nil
	}
	go cli.Connect()

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
