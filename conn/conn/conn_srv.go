package conn

import (
	"net"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
)

const (
	CONN_SRV_CLASS_TCP  = "tcp"
	CONN_SRV_CLASS_UDP  = "udp"
	CONN_SRV_CLASS_QUIC = "quic"

	CONN_SRV_STATE_ON  = 0
	CONN_SRV_STATE_OFF = 1
)

type ConnSrv struct {
	ConnCliSet map[string]*ConnCli
	Port       int
	Class      string

	State int

	listen net.Listener
}

func NewConnSrv(port int, class string) *ConnSrv {
	return &ConnSrv{
		ConnCliSet: make(map[string]*ConnCli),
		Port:       port,
		Class:      class,
	}
}

func (srv *ConnSrv) Run() int {
	srv.State = CONN_SRV_STATE_ON

	res := 0
	switch srv.Class {
	case CONN_SRV_CLASS_TCP:
		res = srv.runTcp()
	case CONN_SRV_CLASS_UDP:
		res = srv.runUdp()
	case CONN_SRV_CLASS_QUIC:
		res = srv.runQuic()
	}

	if res != 0 {
		return -1
	} else {
		//TODO: task listen implement
		go srv.taskListen()
		go srv.taskConnClean()
	}

	return 0
}

func (srv *ConnSrv) runTcp() int {
	listen, _ := net.Listen("tcp", ":"+string(srv.Port))
	for srv.State == CONN_SRV_STATE_ON {
		conn, err := listen.Accept()
		if err != nil {
		} else {
			if len(srv.ConnCliSet) > 1000 {
				ulog.Log().I("conn_srv", "max connection reached, closing")
				conn.Close()
			} else {
				go srv.taskConn(conn)
			}

		}
	}
	return 0
}

func (srv *ConnSrv) runUdp() int {
	laddr, err := net.ResolveUDPAddr("udp", ":"+string(srv.Port))
	if err != nil {
		return -1
	}

	for srv.State == CONN_SRV_STATE_ON {
		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
		} else {
			go srv.taskConn(conn)
		}
	}

	return 0
}

func (srv *ConnSrv) runQuic() int {
	return 0
}

// TODO:
func (srv *ConnSrv) taskListen() {

}

func (srv *ConnSrv) taskConnClean() {
	defer srv.listen.Close()
	tick := time.NewTicker(2 * time.Second)
	for range tick.C {
		if srv.State == CONN_SRV_STATE_OFF {
			break
		}

		for _, cli := range srv.ConnCliSet {
			if cli.State == CONN_CLI_STATE_DISCONNECTED {
				cli.Shutdown()
				delete(srv.ConnCliSet, cli.Conn.RemoteAddr())
			}

		}
	}
}

/**
 * handling messages in flow
 */
func (srv *ConnSrv) taskConn(conn net.Conn) {
	// coder := &Coder{}

	// cli := NewConnCli(conn, coder)
	// srv.AddConnCli(cli)
	// cli.Run()

}

func (c *ConnSrv) AddConnCli(cli ConnCli) {
	// c.ConnCliSet[cli.conn.RemoteAddr()] = cli
}
