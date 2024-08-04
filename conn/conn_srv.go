package conn

import (
	"time"

	"github.com/lingfliu/ucs_core/coder"
)

const (
	STATE_ON  = 0
	STATE_OFF = 1
)

type Srv struct {
	ConnCfg  *ConnCfg
	Codebook *coder.Codebook
	CliSet   map[string]*Cli
	State    int
	Conn     Conn
}

func (srv *Srv) Start() {
	srv.State = STATE_ON
	switch srv.ConnCfg.Class {
	case CONN_CLASS_TCP:
		srv.Conn = NewTcpConn(srv.ConnCfg)
	case CONN_CLASS_UDP:
		srv.Conn = NewUdpConn(srv.ConnCfg)
	default:
		srv.Conn = NewTcpConn(srv.ConnCfg)
	}

	chn := make(chan Conn)
	srv.Conn.Listen(chn)

	go srv._task_accept(chn)

	go srv._task_cleanup()
}

func (srv *Srv) Stop() {
	srv.State = STATE_OFF
	for _, v := range srv.CliSet {
		v.Close()
	}
}

func (srv *Srv) _task_accept(chn chan Conn) {
	for srv.State == STATE_ON {
		select {
		case c := <-chn:
			cli := &Cli{
				Conn:     c,
				Codebook: srv.Codebook,
				Coder:    coder.NewUCoderFromCodebook(srv.Codebook),
				State:    CLI_STATE_CONNECTED,
			}

			srv.CliSet[cli.Conn.GetRemoteAddr()] = cli
			//TODO: cli working
		}
	}
}

func (srv *Srv) _task_cleanup() {
	tic := time.NewTicker(time.Second * 1)
	for srv.State == STATE_ON {
		select {
		case <-tic.C:
			// for _, v := range srv.CliSet {
			// }
		}
	}
}
