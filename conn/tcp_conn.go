package conn

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

type TcpConn struct {
	BaseConn
	c         *net.TCPConn
	stateLock *sync.RWMutex
}

func NewTcpConn(cfg *ConnCfg) *TcpConn {
	sigRun, cancelRun := context.WithCancel(context.Background())
	sigRw, cancelRw := context.WithCancel(context.Background())
	c := &TcpConn{
		BaseConn: BaseConn{
			State:          CONN_STATE_DISCONNECTED,
			RemoteAddr:     cfg.RemoteAddr,
			Port:           cfg.Port,
			Class:          CONN_CLASS_TCP,
			KeepAlive:      cfg.KeepAlive,
			ReconnectAfter: cfg.ReconnectAfter,
			Timeout:        cfg.Timeout,
			TimeoutRw:      cfg.TimeoutRw,
			Rx:             make(chan []byte, 32),
			Tx:             make(chan []byte, 32),
			Io:             make(chan []chan []byte),

			sigRun:    sigRun,
			cancelRun: cancelRun,
			sigRw:     sigRw,
			cancelRw:  cancelRw,
		},
		c:         nil,
		stateLock: &sync.RWMutex{},
	}
	return c
}

/*
 * Connect to remote addr
 * @return 0 if connected, -1 if failed, -2 if already connected
 */
func (c *TcpConn) Connect() int {

	state := c.State
	if state == CONN_STATE_CONNECTED || state == CONN_STATE_CONNECTING {
		return -2
	}

	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	tcp, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.RemoteAddr, c.Port), time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		ulog.Log().I("tcpconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED

		return -1
	} else {
		c.c = tcp.(*net.TCPConn)
		// c.Rx = make(chan []byte, 32)
		// c.Tx = make(chan []byte, 32)

		// c.Io <- []chan []byte{c.Rx, c.Tx}

		go c._task_recv(c.sigRun)
		go c._task_send(c.sigRun)
		return 0
	}
}

/*
 * Disconnect from remote addr
 * return 0 anyway
 */
func (c *TcpConn) Disconnect() int {
	if c.State == CONN_STATE_DISCONNECTED {
		return 0
	}
	c.lastDisconnectAt = utils.CurrentTime()
	c.State = CONN_STATE_DISCONNECTED

	close(c.Rx)
	c.Io <- []chan []byte{c.Rx, c.Tx}

	if c.c != nil {
		c.c.Close()
	}
	return 0
}

/*
 * Listen on port
 * @param ch channel to send new conn
 * @return 0 if success, -1 if failed
 */
func (c *TcpConn) Listen(ctx context.Context, ch chan Conn) {
	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: c.Port,
	}
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		ulog.Log().E("tcpconn", "listen failed, check port")
		c.State = CONN_STATE_CLOSED
		c.Close()
		return
	}

	c.State = CONN_STATE_LISTENING
	for c.State != CONN_STATE_CLOSED {
		select {
		case <-ctx.Done():
			return
		default:
			cc, err := l.Accept()
			if err != nil {
				ulog.Log().E("tcpconn", "accept failed, skip")
				continue
			}
			// create tcpconn using default conn cfg
			cfg := ctx.Value("conn_cfg").(*ConnCfg)
			tcp := NewTcpConn(cfg)
			tcp.c = cc.(*net.TCPConn)
			tcp.Tx = make(chan []byte, 32)
			tcp.Rx = make(chan []byte, 32)
			ch <- tcp //TODO: test the channel for new conn handling
		}
	}
}

func (c *TcpConn) Close() int {
	if c.State == CONN_STATE_CLOSED {
		return -2
	}

	c.State = CONN_STATE_CLOSED
	if c.c != nil {
		c.c.Close()
	}
	c.Io <- make([]chan []byte, 0)
	return 0
}

func (c *TcpConn) Read(bs []byte) int {
	if c.State != CONN_STATE_CONNECTED {
		return -1
	}

	c.c.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
	n, err := c.c.Read(bs)
	if err != nil {
		ulog.Log().E("tcpconn", "read failed")
		c.Disconnect()
		return -1
	}
	return n
}

func (c *TcpConn) Write(bs []byte) int {
	if c.State != CONN_STATE_CONNECTED {
		return -2
	}

	c.c.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
	n, err := c.c.Write(bs)
	if err != nil {
		ulog.Log().E("tcpconn", "write failed")
		c.Disconnect()
		return -1
	}
	return n
}

func (c *TcpConn) GetRemoteAddr() string {
	return c.RemoteAddr
}

func (c *TcpConn) GetState() int {
	return c.State
}

func (c *TcpConn) GetRx() chan []byte {
	return c.Rx
}

func (c *TcpConn) GetTx() chan []byte {
	return c.Tx
}

func (c *TcpConn) Start(sigRun context.Context) chan []chan []byte {
	go c._task_connect(sigRun)
	return c.Io
}

//tasks

func (c *TcpConn) _task_recv(ctx context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case <-ctx.Done():
			return
		default:
			bs := make([]byte, 1024)
			c.Read(bs)
			c.Rx <- bs
		}
	}
}

func (c *TcpConn) _task_send(ctx context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case buff := <-c.Tx:
			if len(buff) == 0 {
				continue
			} else {
				n := c.Write(buff)
				if n == -1 {
					ulog.Log().E("tcpconn", "write failed, disconnect")
					c.Disconnect()
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *TcpConn) _task_connect(ctx context.Context) {
	if !c.KeepAlive {
		return
	}

	tic := time.NewTicker(time.Duration(1) * time.Second)

	for c.State != CONN_STATE_CLOSED {
		select {
		case <-tic.C:
			if c.State == CONN_STATE_DISCONNECTED && utils.CurrentTime()-c.lastDisconnectAt > c.ReconnectAfter {
				c.Connect()
			}
		case <-ctx.Done():
			return
		}
	}
}
