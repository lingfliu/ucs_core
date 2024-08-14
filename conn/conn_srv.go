package conn

import (
	"context"
	"strconv"
	"time"

	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

const (
	SRV_STATE_ON  = 0
	SRV_STATE_OFF = 1
)

type ConnSrv struct {
	ConnCfg  *ConnCfg
	Codebook *coder.Codebook
	CliSet   map[string]*ConnCli
	State    int
	Conn     Conn

	//cli cfg
	MsgTimeout int64

	//sigs
	sigRun    context.Context
	cancelRun context.CancelFunc

	ctxCfg context.Context

	//callbacks
	MsgHandler func(*ConnCli, *coder.ZeroMsg)
}

func NewConnSrv(connCfg *ConnCfg, cb *coder.Codebook) *ConnSrv {
	sigRun, cancelRun := context.WithCancel(context.Background())
	cfg := context.WithValue(context.Background(), utils.CtxKeyCfg{}, connCfg)
	srv := &ConnSrv{
		ConnCfg:  connCfg,
		Codebook: cb,
		CliSet:   make(map[string]*ConnCli),
		State:    SRV_STATE_OFF,
		Conn:     nil,

		sigRun:    sigRun,
		cancelRun: cancelRun,

		ctxCfg: cfg,
	}
	return srv
}

func (srv *ConnSrv) Start() {
	srv.State = SRV_STATE_ON
	switch srv.ConnCfg.Class {
	case CONN_CLASS_TCP:
		srv.Conn = NewTcpConn(srv.ConnCfg)
	case CONN_CLASS_UDP:
		// srv.Conn = NewUdpConn(srv.ConnCfg)
	default:
		srv.Conn = NewTcpConn(srv.ConnCfg)
	}

	connChn := make(chan Conn)
	go srv.Conn.Listen(srv.sigRun, srv.ctxCfg, connChn)

	go srv._task_spawn_cli(connChn)
	go srv._task_cleanup()
}

func (srv *ConnSrv) Stop() {
	srv.State = SRV_STATE_OFF
	srv.cancelRun()
	for _, v := range srv.CliSet {
		v.Close()
	}
	for k := range srv.CliSet {
		delete(srv.CliSet, k)
	}
	srv.Conn.Close()
}

func (srv *ConnSrv) SpawnConnCli(c Conn) *ConnCli {
	if c == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	codebook := srv.Codebook

	cli := &ConnCli{
		State: CLI_STATE_CONNECTED,
		Conn:  c,
		Mode:  CLI_MODE_SPAWN,

		Codebook: codebook,
		Coder:    coder.NewZeroCoder(),

		txMsg: make(chan *coder.ZeroMsg, 32),

		sigCoding:    ctx,
		cancelCoding: cancel,

		sigRun:    ctx2,
		cancelRun: cancel2,

		lastMsgAt: utils.CurrentTime(),
	}
	return cli
}

func (srv *ConnSrv) _task_spawn_cli(connChn chan Conn) {
	for srv.State == SRV_STATE_ON {
		select {
		case c := <-connChn:
			cli := srv.SpawnConnCli(c)
			ulog.Log().I("conncli", "new cli spawned "+cli.Conn.GetRemoteAddr())
			srv.CliSet[cli.Conn.GetRemoteAddr()] = cli

			cli.HandleMsg = func(msg *coder.ZeroMsg) {
				if srv.MsgHandler != nil {
					srv.MsgHandler(cli, msg)
				}
			}
			srv.PrepareConn(cli)
			cli.StartRw()
		case <-srv.sigRun.Done():
			return
		}
	}
}

func (srv *ConnSrv) PrepareConn(cli *ConnCli) {
	//non-blocking msg sending
	// cli.PushTxMsg(cli.Coder.CreatePingpongMsg())
}

func (srv *ConnSrv) _task_cleanup() {
	tic := time.NewTicker(time.Second * 1)
	for srv.State == SRV_STATE_ON {
		select {
		case <-tic.C:
			ulog.Log().I("connsrv", "cleaning up inactive cli")
			for k, cli := range srv.CliSet {
				if utils.CurrentTime()-cli.lastMsgAt > srv.MsgTimeout {
					//inactive cli, remove
					ulog.Log().I("connsrv", "removing cli: "+k+" inactive for: "+strconv.FormatInt((utils.CurrentTime()-cli.lastMsgAt)/1000/1000, 10)+" ms")
					cli.Close()
					delete(srv.CliSet, k)
				}
			}
		case <-srv.sigRun.Done():
			return
		}
	}
}
