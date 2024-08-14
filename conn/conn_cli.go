package conn

import (
	"context"
	"time"

	"github.com/lingfliu/ucs_core/coder"
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

	lastMsgAt int64

	//callbacks
	//TODO: add internal handling
	HandleMsg func(*coder.ZeroMsg)
}

func NewConnCli(connCfg *ConnCfg, cb *coder.Codebook) *ConnCli {
	ctx, cancel := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	codebook := cb
	var c Conn
	switch connCfg.Class {
	case CONN_CLASS_TCP:
		c = NewTcpConn(connCfg)
	// case CONN_CLASS_UDP:
	// c = NewUdpConn(connCfg)
	// case CONN_CLASS_QUIC:
	// c = NewQuicConn(connCfg)
	default:
		c = NewTcpConn(connCfg)
	}

	if c == nil {
		return nil
	}

	cli := &ConnCli{
		State: CLI_STATE_DISCONNECTED,
		Conn:  c,
		Mode:  CLI_MODE_HOST,

		Codebook: codebook,
		Coder:    coder.NewZeroCoder(),

		txMsg: make(chan *coder.ZeroMsg, 32),

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

func (cli *ConnCli) _task_handle_connect(io chan int) {
	for cli.State != CLI_STATE_CLOSED {
		select {
		case cs := <-io:
			if cs == CONN_STATE_CLOSED {
				//conn closed
				cli.cancelCoding()
			} else {

				//TODO: finish previous decode & encode routines
				cli.cancelCoding()

				ctx, cancel := context.WithCancel(context.Background())
				cli.sigCoding = ctx
				cli.cancelCoding = cancel

				//restart decode / encode routines
				cli.State = CLI_STATE_CONNECTED
				cli.StartRw()
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

	if cli.Mode == CLI_MODE_SPAWN {
		go cli._task_pingpong()
	}
}

func (cli *ConnCli) Disconnect() {
	cli.Conn.Disconnect()
}

func (cli *ConnCli) Close() {
	cli.State = CLI_STATE_CLOSED
	cli.Conn.Close()
	cli.cancelRun()
}

func (cli *ConnCli) PushTxMsg(msg *coder.ZeroMsg) {
	cli.txMsg <- msg
}
