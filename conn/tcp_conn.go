package conn

import (
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
			TxBuff:         utils.NewByteArrayRingBuffer(32, 1024),
		},
	}
	return c
}

func (c *TcpConn) Connect() {
	if c.State == CONN_STATE_CONNECTED || c.State == CONN_STATE_CONNECTING {
		return
	}

	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	tcp, err := net.DialTimeout("tcp", c.RemoteAddr, time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		ulog.Log().I("tcpconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED
	} else {
		c.State = CONN_STATE_CONNECTED
		c.c = tcp.(*net.TCPConn)

		go c._task_recv()
		go c._task_send()
	}
}

func (c *TcpConn) Disconnect() {
	if c.State == CONN_STATE_DISCONNECTED {
		return
	}

	c.lastDisconnectAt = utils.CurrentTime()
	c.State = CONN_STATE_DISCONNECTED
	if c.c != nil {
		c.c.Close()
	}
}

func (c *TcpConn) Listen(ch chan Conn) {
	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: c.Port,
	}
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		ulog.Log().E("tcpconn", "listen failed, check port")
		c.Close()
	}

	for c.State != CONN_STATE_CLOSE {
		cc, err := l.Accept()
		if err != nil {
			ulog.Log().E("tcpconn", "accept failed, shutdown")
			break
		}
		tcp := &TcpConn{
			BaseConn: BaseConn{
				State:      CONN_STATE_CONNECTED,
				Class:      CONN_CLASS_TCP,
				KeepAlive:  c.KeepAlive,
				Timeout:    c.Timeout,
				TimeoutRw:  c.TimeoutRw,
				LocalAddr:  c.LocalAddr,
				RemoteAddr: cc.RemoteAddr().String(),
				Port:       c.Port,
				TxBuff:     utils.NewByteArrayRingBuffer(32, 1024),
			},
			c: cc.(*net.TCPConn),
		}
		ch <- tcp //TODO: test the channel for new conn handling
	}
}

func (c *TcpConn) _task_recv() {
	buff := make([]byte, 1024)
	for {
		c.c.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
		n, err := c.c.Read(buff)

		if err != nil {
			//TODO: handling disconnect
			c.Disconnect()
		}
		if n > 0 {
			c.RxBuff.Push(buff, n)
		}
	}
}

func (c *TcpConn) StartRecv() {
	go c._task_recv()
}

func (c *TcpConn) Read(buff []byte) int {
	c.c.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
	n, err := c.c.Read(buff)
	if err != nil {
		//read error, break the connection
		c.Disconnect()
		return -1
	} else {
		return n
	}
}

func (c *TcpConn) InstantWrite(buff []byte) int {
	err := c.c.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
	if err != nil {
		c.Disconnect()
	}

	n, err := c.c.Write(buff)
	if err != nil {
		//write error, break the connection
		c.Disconnect()
		return -1
	}
	return n
}

func (c *TcpConn) ScheduleWrite(buff []byte) {
	c.TxBuff.Push(buff)
}

func (c *TcpConn) _task_send() {
	for c.State == CONN_STATE_CONNECTED {
		buff := c.TxBuff.Pop()
		if buff == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		err := c.c.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
		if err != nil {
			c.Disconnect()
		}

		_, err = c.c.Write(buff)
		if err != nil {
			c.Disconnect()
		}
	}
}

func (c *TcpConn) _task_connect() {
	tic := time.NewTicker(time.Second * 1)
	for c.State != CONN_STATE_CLOSE {
		select {
		case <-tic.C:
			if c.State == CONN_STATE_DISCONNECTED && utils.CurrentTime()-c.lastDisconnectAt > c.ReconnectAfter {
				c.Connect()
			}
		}
	}
}

func (c *TcpConn) Close() {
	c.State = CONN_STATE_CLOSE
	if c.c != nil {
		c.c.Close()
	}
}

func (c *TcpConn) GetRemoteAddr() string {
	return c.RemoteAddr
}
