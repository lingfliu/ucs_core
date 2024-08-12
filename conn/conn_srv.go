package conn

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

const (
	STATE_ON  = 0
	STATE_OFF = 1
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

	//callbacks
	MsgHandler func(*ConnCli, *coder.ZeroMsg)
}

func NewConnSrv(connCfg *ConnCfg, cb *coder.Codebook) *ConnSrv {
	sigRun, cancelRun := context.WithCancel(context.Background())
	srv := &ConnSrv{
		ConnCfg:  connCfg,
		Codebook: cb,
		CliSet:   make(map[string]*ConnCli),
		State:    STATE_OFF,
		Conn:     nil,

		sigRun:    sigRun,
		cancelRun: cancelRun,
	}
	return srv
}

func (srv *ConnSrv) Start() {
	srv.State = STATE_ON
	switch srv.ConnCfg.Class {
	case CONN_CLASS_TCP:
		srv.Conn = NewTcpConn(srv.ConnCfg)
	case CONN_CLASS_UDP:
		// srv.Conn = NewUdpConn(srv.ConnCfg)
	default:
		srv.Conn = NewTcpConn(srv.ConnCfg)
	}

	connChn := make(chan Conn)
	srv.Conn.Listen(srv.sigRun, connChn)

	go srv._task_spawn_cli(connChn)
	go srv._task_cleanup()
}

func (srv *ConnSrv) Stop() {
	srv.State = STATE_OFF
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
	}
	return cli
}

func (srv *ConnSrv) _task_spawn_cli(connChn chan Conn) {
	for srv.State == STATE_ON {
		select {
		case c := <-connChn:
			cli := srv.SpawnConnCli(c)
			srv.CliSet[cli.Conn.GetRemoteAddr()] = cli

			cli.HandleMsg = func(msg *coder.ZeroMsg) {

				//TODO: demo code, remove on release
				ulog.Log().I("conncli", "received msg")
				jsonStr, err := json.Marshal(msg)
				if err != nil {
					ulog.Log().E("conncli", "msg decode error")
				}
				ulog.Log().I("conncli", "received msg: "+string(jsonStr))

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
	cli.PushTxMsg(cli.Coder.CreatePingpongMsg())
}

func (srv *ConnSrv) _task_cleanup() {
	tic := time.NewTicker(time.Second * 1)
	for srv.State == STATE_ON {
		select {
		case <-tic.C:
			for k, cli := range srv.CliSet {
				if utils.CurrentTime()-cli.lastMsgAt > srv.MsgTimeout {
					//inactive cli, remove
					delete(srv.CliSet, k)
				}
			}
		case <-srv.sigRun.Done():
			return
		}
	}
}
