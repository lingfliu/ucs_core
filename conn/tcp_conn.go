package conn

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

type TcpConn struct {
	BaseConn
	c *net.TCPConn
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
			Io:             make(chan int),

			sigRun:    sigRun,
			cancelRun: cancelRun,
			sigRw:     sigRw,
			cancelRw:  cancelRw,

			lastDisconnectAt: 0,
			lastConnectAt:    0,
			lastRecvAt:       0,
		},
		c: nil,
	}
	return c
}

/*
 * Connect to remote addr
 * @return 0 if connected, -1 if failed, -2 if already connected
 */
func (c *TcpConn) Connect() int {

	state := c.State
	if state == CONN_STATE_CONNECTED {
		return -2
	}
	if state == CONN_STATE_CONNECTING {
		return -2
	}

	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	tcp, err := net.DialTimeout("tcp", utils.IpPortJoin(c.RemoteAddr, c.Port), time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		ulog.Log().I("tcpconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED

		return -1
	} else {
		//renew conn and rw control
		c.State = CONN_STATE_CONNECTED
		c.Io <- CONN_STATE_CONNECTED
		c.c = tcp.(*net.TCPConn)
		c.sigRw, c.cancelRw = context.WithCancel(context.Background())
		go c._task_recv(c.sigRw)
		go c._task_send(c.sigRw)
		return 0
	}
}

/*
 * Disconnect from remote addr
 * return 0 anyway
 */
func (c *TcpConn) Disconnect() int {
	if c.State == CONN_STATE_DISCONNECTED {
		return -2
	}

	c.lastDisconnectAt = utils.CurrentTime()
	c.State = CONN_STATE_DISCONNECTED

	c.Io <- c.State

	c.cancelRw()
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
func (c *TcpConn) Listen(sigRun context.Context, ctxCfg context.Context, ch chan Conn) {
	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: c.Port,
	}
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		ulog.Log().E("tcpconn", "listen failed, check port")
		c.Close()
		return
	}

	c.State = CONN_STATE_LISTENING
	for c.State != CONN_STATE_CLOSED {
		select {
		case <-sigRun.Done():
			return
		default:
			cc, err := l.Accept()
			if err != nil {
				ulog.Log().E("tcpconn", "accept failed, skip")
				continue
			}
			// create tcpconn using default conn cfg
			cfg := ctxCfg.Value(utils.CtxKeyCfg{}).(*ConnCfg)
			tcp := NewTcpConn(cfg)
			tcp.RemoteAddr = cc.RemoteAddr().String()

			tcp.c = cc.(*net.TCPConn)
			tcp.State = CONN_STATE_CONNECTED
			tcp.lastConnectAt = utils.CurrentTime()
			tcp.sigRun, tcp.cancelRun = context.WithCancel(context.Background())
			tcp.sigRw, tcp.cancelRw = context.WithCancel(context.Background())
			go tcp._task_recv(tcp.sigRw)
			go tcp._task_send(tcp.sigRw)

			ch <- tcp
		}
	}
}
func (c *TcpConn) Start(sigRun context.Context) chan int {
	go c._task_connect(sigRun)
	return c.Io
}

func (c *TcpConn) Close() int {
	if c.State == CONN_STATE_CLOSED {
		return -2
	}

	c.State = CONN_STATE_CLOSED
	c.Io <- CONN_STATE_CLOSED
	c.cancelRw()
	c.cancelRun()
	if c.c != nil {
		c.c.Close()
	}
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

// tasks
func (c *TcpConn) _task_recv(sigRw context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case <-sigRw.Done():
			return
		default:
			bs := make([]byte, 1024)
			n := c.Read(bs)
			if n > 0 {
				c.Rx <- bs[:n]
			}
		}
	}
}

func (c *TcpConn) _task_send(sigRw context.Context) {
	ulog.Log().I("tcpconn", "task send started")
	for c.State == CONN_STATE_CONNECTED {
		select {
		case buff := <-c.Tx:
			if len(buff) == 0 {
				continue
			} else {
				c.Write(buff)
			}
		case <-sigRw.Done():
			return
		}
	}
}

func (c *TcpConn) _task_connect(sigRun context.Context) {
	if !c.KeepAlive {
		return
	}

	//first connect
	c.Connect()

	tic := time.NewTicker(time.Duration(1) * time.Second)

	for c.State != CONN_STATE_CLOSED {
		select {
		case <-tic.C:
			if c.State == CONN_STATE_DISCONNECTED && utils.CurrentTime()-c.lastDisconnectAt > c.ReconnectAfter {
				c.Connect()
			}
		case <-sigRun.Done():
			return
		}
	}
}
