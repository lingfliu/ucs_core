package conn

import (
	"context"
	"time"

	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

const (
	CLI_STATE_CLOSED       = 0
	CLI_STATE_DISCONNECTED = 1
	CLI_STATE_CONNECTING   = 2
	CLI_STATE_CONNECTED    = 3
	CLI_STATE_AUTH         = 4

	CLI_MODE_HOST  = 1 // cli side cli
	CLI_MODE_SPAWN = 2 // server side cli
)

type ConnCli struct {
	State int
	Conn  Conn

	Codebook *coder.Codebook
	Coder    *coder.ZeroCoder

	Token string //ascii string token

	Mode int

	txMsg chan *coder.ZeroMsg
	//event
	OnReq func(*coder.ZeroMsg) *coder.ZeroMsg

	//sigs
	sigCoding    context.Context
	cancelCoding context.CancelFunc

	sigRun    context.Context
	cancelRun context.CancelFunc

	lastMsgAt  int64
	MsgTimeout int64

	//callbacks
	//TODO: add internal handling
	HandleMsg func(*coder.ZeroMsg)
}

func NewConnCli(connCfg *ConnCfg, cb *coder.Codebook) *ConnCli {
	codebook := cb
	var c Conn
	switch connCfg.Class {
	case CONN_CLASS_TCP:
		c = NewTcpConn(connCfg)
	case CONN_CLASS_UDP:
		c = NewUdpConn(connCfg)
	case CONN_CLASS_QUIC:
		c = NewQuicConn(connCfg)
	default:
		c = NewTcpConn(connCfg)
	}

	if c == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	cli := &ConnCli{
		State: CLI_STATE_DISCONNECTED,
		Conn:  c,
		Mode:  CLI_MODE_HOST,

		Codebook: codebook,
		Coder:    coder.NewZeroCoder(),

		txMsg:      make(chan *coder.ZeroMsg, 32),
		MsgTimeout: 5 * 1000 * 1000 * 1000,

		sigCoding:    ctx,
		cancelCoding: cancel,
		sigRun:       ctx2,
		cancelRun:    cancel2,
	}
	return cli
}

func (cli *ConnCli) _task_decode(rx chan []byte) {
	for cli.State == CLI_STATE_CONNECTED || cli.State == CLI_STATE_AUTH {
		select {
		case bs := <-rx:
			if len(bs) < 1 {
				continue
			}
			cli.Coder.FastDecode(bs)
		case <-cli.sigCoding.Done():
			return
		}
	}
}

func (cli *ConnCli) _task_handle_msg() {
	for cli.State == CLI_STATE_CONNECTED || cli.State == CLI_STATE_AUTH {
		select {
		case msg := <-cli.Coder.RxMsg:
			cli.lastMsgAt = utils.CurrentTime()
			if cli.HandleMsg != nil {
				cli.HandleMsg(msg)
			}
		case <-cli.sigRun.Done():
			return
		}
	}
}

func (cli *ConnCli) _task_pingpong() {
	tic := time.NewTicker(time.Duration(1) * time.Second)
	for cli.State == CLI_STATE_CONNECTED || cli.State == CLI_STATE_AUTH {
		select {
		case <-tic.C:
			cli.txMsg <- &coder.ZeroMsg{
				Class: 0,
			}
		case <-cli.sigCoding.Done():
			return
		}
	}
}

func (cli *ConnCli) _task_encode(tx chan []byte) {
	for cli.State == CLI_STATE_CONNECTED || cli.State == CLI_STATE_AUTH {
		select {
		case msg := <-cli.txMsg:
			bs := make([]byte, 1024)
			n := cli.Coder.Encode(msg, bs)
			tx <- bs[:n]
		case <-cli.sigCoding.Done():
			return
		}
	}
}

func (cli *ConnCli) _task_keepalive(io chan int) {
	tic := time.NewTicker(time.Duration(1) * time.Second)
	for cli.State == CLI_STATE_CONNECTED {
		select {
		case <-tic.C:
			if utils.CurrentTime()-cli.lastMsgAt > cli.MsgTimeout {
				ulog.Log().E("conncli", "msg timeout")
				cli.Disconnect()
				return
			}
		case <-cli.sigRun.Done():
			return
		case state := <-io:
			if state == CONN_STATE_DISCONNECTED {
				return
			}
		}
	}
}

func (cli *ConnCli) _task_handle_connect(io chan int) {
	for cli.State != CLI_STATE_CLOSED {
		select {
		case cs := <-io:
			switch cs {
			case CONN_STATE_CLOSED:
				ulog.Log().I("conncli", "conn closed")
				//conn closed
				cli.cancelCoding()
			case CONN_STATE_DISCONNECTED:
				cli.State = CLI_STATE_DISCONNECTED
			case CONN_STATE_CONNECTING:
				cli.State = CLI_STATE_CONNECTING
			case CONN_STATE_CONNECTED:
				cli.State = CLI_STATE_CONNECTED
				cli.lastMsgAt = utils.CurrentTime()
				cli.cancelCoding()
				//TODO: finish previous decode & encode routines
				ctx, cancel := context.WithCancel(context.Background())
				cli.sigCoding = ctx
				cli.cancelCoding = cancel
				//restart decode / encode routines
				cli.StartRw()

				go cli._task_keepalive(io)
			}

		case <-cli.sigRun.Done():
			cli.cancelCoding()
			return
		}
	}
}

/**
 * effective under spawn mode
 **/
func (cli *ConnCli) Start() {
	sigRun := cli.sigRun
	io := cli.Conn.Start(sigRun)
	go cli._task_handle_connect(io)
}

func (cli *ConnCli) StartRw() {
	rx := cli.Conn.GetRx()
	tx := cli.Conn.GetTx()

	go cli._task_decode(rx)
	go cli._task_encode(tx)
	go cli._task_handle_msg()
	go cli._task_handle_connect(cli.Conn.GetIo())
	go cli._task_pingpong()
}

func (cli *ConnCli) Disconnect() {
	cli.Conn.Disconnect()
}

func (cli *ConnCli) Stop() {
	cli.State = CLI_STATE_CLOSED
	cli.Conn.Close()
	cli.cancelRun()
}

func (cli *ConnCli) PushTxMsg(msg *coder.ZeroMsg) {
	cli.txMsg <- msg
}
