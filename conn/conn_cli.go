package conn

import (
	"context"

	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/ulog"
)

const (
	CLI_STATE_CLOSED       = 0
	CLI_STATE_DISCONNECTED = 1
	CLI_STATE_CONNECTING   = 2
	CLI_STATE_CONNECTED    = 3
	CLI_STATE_AUTH         = 4

	CLI_MODE_HOST  = 1 // cli at cli side
	CLI_MODE_SPAWN = 2 // cli at server side
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
			break
		case <-cli.sigCoding.Done():
			return
		}
	}
}

func (cli *ConnCli) _task_handle_msg() {
	for cli.State == CLI_STATE_CONNECTED || cli.State == CLI_STATE_AUTH {
		select {
		case msg := <-cli.Coder.RxMsg:
			if cli.HandleMsg != nil {
				cli.HandleMsg(msg)
			}
			break
		case <-cli.sigRun.Done():
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

func (cli *ConnCli) _task_handle_connect(io chan []chan []byte) {
	for cli.State == STATE_ON {
		select {
		case cs := <-io:
			if len(cs) == 0 {
				//conn closed
				cli.cancelCoding()
			} else {
				rx := cs[0]
				tx := cs[1]

				//TODO: finish previous decode & encode routines
				cli.cancelCoding()

				ctx, cancel := context.WithCancel(context.Background())
				cli.sigCoding = ctx
				cli.cancelCoding = cancel

				//restart decode / encode routines
				go cli._task_decode(rx)
				go cli._task_encode(tx)
			}
		case <-cli.sigRun.Done():
			cli.cancelCoding()
			return
		}
	}
}

func (cli *ConnCli) Connect() {
	cli.State = CLI_STATE_CONNECTING
	ret := cli.Conn.Connect()
	if ret != 0 {
		cli.State = CLI_STATE_DISCONNECTED
		ulog.Log().E("ConnCli", "Connect failed")
	} else {
		cli.State = CLI_STATE_CONNECTED
		cli.StartRw()
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
	io := cli.Conn.GetIo()

	go cli._task_decode(rx)
	go cli._task_encode(tx)
	go cli._task_handle_connect(io)
	go cli._task_handle_msg()
}

func (cli *ConnCli) Disconnect() {
	cli.Conn.Disconnect()
}

func (cli *ConnCli) Close() {
	cli.State = CLI_STATE_CLOSED
	cli.Conn.Close()
	cli.cancelRun()

	close(cli.Conn.GetTx())
}

func (cli *ConnCli) PushTxMsg(msg *coder.ZeroMsg) {
	cli.txMsg <- msg
}
