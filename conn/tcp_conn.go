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

	OnConnected    func()
	OnDisconnected func()

	RxBuff *utils.ByteRingBuffer
	TxBuff *utils.ByteArrayRingBuffer
}

func NewTcpConn(remoteAddr string, port int) *TcpConn {
	c := &TcpConn{
		BaseConn: BaseConn{
			State:      CONN_STATE_DISCONNECTED,
			RemoteAddr: remoteAddr,
			Port:       port,
			Class:      CONN_CLASS_TCP,
		},
		RxBuff: utils.NewByteRingBuffer(1024),
		TxBuff: utils.NewByteArrayRingBuffer(32, 1024),
	}
	return c
}

func (c *TcpConn) Connect() int {
	c.State = CONN_STATE_CONNECTING
	tcp, err := net.DialTimeout("tcp", c.RemoteAddr, time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		ulog.Log().I("tcpconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED
		return -1
	}

	c.State = CONN_STATE_CONNECTED
	c.c = tcp.(*net.TCPConn)

	go c._task_recv()
	go c._task_send()
	return 0
}

func (c *TcpConn) Establish(tcp *net.TCPConn) int {
	c.State = CONN_STATE_CONNECTED
	c.c = tcp

	go c._task_recv()
	go c._task_send()
	return 0
}

func (c *TcpConn) Disconnect() int {
	if c.c != nil {
		c.c.Close()
		//TODO: finish rx & tx
		c.State = CONN_STATE_DISCONNECTED
	}
	return 0
}

func (c *TcpConn) _task_recv() {
	buff := make([]byte, 1024)
	for {
		c.c.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
		n, err := c.c.Read(buff)

		if err != nil {
			//TODO: handling disconnect
			// return -1
			c.Disconnect()
		}
		if n > 0 {
			c.RxBuff.Push(buff, n)
		}
	}
}

func (c *TcpConn) Read(buff []byte) int {
	return 0
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
			time.Sleep(10 * time.Millisecond)
			continue
		}

		err := c.c.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
		if err != nil {
			c.Disconnect()
		}

		_, err = c.c.Write(buff)
		if err != nil {
			//TODO: handle write error
			c.Disconnect()
		}

		//TODO: handle write success
		// return n
	}
}
