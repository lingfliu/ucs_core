package conn

import (
	"net"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
)

const (
	STATE_START = 0
	STATE_STOP  = 1
)

type TcpConnSrv struct {
	Port          int
	TcpConnCliSet map[string]*TcpConnCli
	State         int
}

func (srv *TcpConnSrv) Start() int {
	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: srv.Port,
	}
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		ulog.Log().E("tcpconn", "listen failed, check port")
		return -1
	}
	go srv._task_listen(l)
	return 0
}

func (c *TcpConnSrv) _task_listen(listener *net.TCPListener) {
	for {
		cli, err := listener.Accept()
		if err != nil {
			ulog.Log().E("tcpconn", "accept failed, shutdown")
			break
		}

		cc := cli.(*net.TCPConn)

		tcp := &TcpConn{
			BaseConn: BaseConn{
				State: CONN_STATE_CONNECTED,
				Class: CONN_CLASS_TCP,
			},
			c: cc,
		}

		tcp.c.Close()
	}
}

func (srv *TcpConnSrv) _task_cleanup() {
	tic := time.NewTicker(time.Second * 1)
	for srv.State == STATE_START {
		select {
			case <-tic.C:
				for _, v := range srv.TcpConnCliSet {
				}
	}
}
