package conn

import (
	"context"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
)

const (
	CONN_STATE_DISCONNECTED = 0
	CONN_STATE_CONNECTING   = 1
	CONN_STATE_CONNECTED    = 2
	CONN_STATE_CLOSED       = 3
	CONN_STATE_LISTENING    = 4

	CONN_CLASS_TCP   string = "tcp"
	CONN_CLASS_UDP   string = "udp"
	CONN_CLASS_QUIC  string = "quic"
	CONN_CLASS_MQTT  string = "mqtt"
	CONN_CLASS_OPCUA string = "opcua"
	CONN_CLASS_HTTP  string = "http"
	CONN_CLASS_RTSP  string = "rtsp"
	// CONN_CLASS_kcp string = "kcp" //TODO: implement kcp
)

// ********************************************************
// common attrs for conn
// ********************************************************
type BaseConn struct {
	State      int
	KeepAlive  bool  // by default true
	Timeout    int64 // connect timeout
	TimeoutRw  int64 // read write timeout
	LocalAddr  string
	RemoteAddr string
	Port       int
	Class      string //tcp, quic, http, mqtt

	ReconnectAfter   int64
	lastRecvAt       int64
	lastConnectAt    int64
	lastDisconnectAt int64

	Rx chan []byte
	Tx chan []byte

	Io chan int

	//contexts
	sigRun    context.Context
	cancelRun context.CancelFunc
	sigRw     context.Context
	cancelRw  context.CancelFunc
}

/**********************tasks**************************/
func (c *BaseConn) _task_connect(ctx context.Context) {

	tic := time.NewTicker(time.Duration(1) * time.Second)

	for c.State != CONN_STATE_CLOSED {
		select {
		case <-tic.C:
			ulog.Log().E("baseconn", "task connect not implemented")
		// 	if c.State == CONN_STATE_DISCONNECTED && utils.CurrentTime()-c.lastDisconnectAt > c.ReconnectAfter {
		// 		c.Connect()
		// 	}
		case <-ctx.Done():
			return
		}
	}
}

func (c *BaseConn) _task_recv(ctx context.Context) {

	for c.State == CONN_STATE_CONNECTED {
		select {
		case <-ctx.Done():
			return
		default:
			ulog.Log().E("baseconn", "task recv not implemented")

		}
	}
}

func (c *BaseConn) _task_send(ctx context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case buff := <-c.Tx:
			buff[0] = 0
			ulog.Log().E("baseconn", "task send not implemented")

		case <-ctx.Done():
			return
		}
	}
}

/**********************Conn interface implement**************************/
func (c *BaseConn) GetRx() chan []byte {
	return c.Rx
}

func (c *BaseConn) GetTx() chan []byte {
	return c.Tx
}

func (c *BaseConn) GetIo() chan int {
	return c.Io
}

func (c *BaseConn) GetRemoteAddr() string {
	return c.RemoteAddr
}

func (c *BaseConn) GetState() int {
	return c.State
}

func (c *BaseConn) Read(buff []byte) int {
	ulog.Log().E("BaseConn", "Read() not implemented")
	return 0
}

func (c *BaseConn) Write(buff []byte) int {
	ulog.Log().E("BaseConn", "Write() not implemented")
	return 0
}

func (c *BaseConn) Start(sigRun context.Context) chan int {
	go c._task_connect(sigRun)
	return c.Io
}

func (c *BaseConn) Connect() int {
	ulog.Log().E("BaseConn", "Connect() not implemented")
	return 0
}

func (c *BaseConn) Disconnect() int {
	ulog.Log().E("BaseConn", "Disconnect() not implemented")
	return 0
}

func (c *BaseConn) Listen(sigRun context.Context, ctxCfg context.Context, ch chan Conn) {
	ulog.Log().E("BaseConn", "Listen() not implemented")
}

func (c *BaseConn) Close() int {
	if c.State == CONN_STATE_CLOSED {
		return -2
	}

	c.cancelRw()  //stop rw
	c.cancelRun() //stop connect
	c.State = CONN_STATE_CLOSED
	c.Disconnect()
	close(c.Rx)
	return 0
}

type Conn interface {
	Start(ctx context.Context) chan int
	Connect() int
	Disconnect() int
	Close() int

	Read(bs []byte) int
	Write(bs []byte) int

	//Srv
	Listen(sig context.Context, cfg context.Context, ch chan Conn)

	//Attr fetch
	GetRemoteAddr() string
	GetState() int
	GetRx() chan []byte
	GetTx() chan []byte
	GetIo() chan int
}

type ConnCfg struct {
	RemoteAddr     string
	LocalAddr      string
	Port           int
	Class          string
	KeepAlive      bool
	ReconnectAfter int64
	Timeout        int64
	TimeoutRw      int64
}
