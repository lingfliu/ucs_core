package conn

import (
	"net"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

type TcpConn struct {
	BaseConn
	c    *net.TCPConn
	Mode int // 0 for client, 1 for server

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
	}
	return c
}

func (c *TcpConn) Connect() int {
	c.State = CONN_STATE_CONNECTING
	tcp, err := net.DialTimeout("tcp", c.RemoteAddr, time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		ulog.Log().I("tcpconn", "connect to %s:%d failed", c.RemoteAddr, c.Port)
		c.State = CONN_STATE_DISCONNECTED
		return -1
	}

	c.State = CONN_STATE_CONNECTED
	c.c = tcp

	c.StartRecv()
	c.StartSend()
	return 0
}

func (c *TcpConn) Disconnect() int {
	if c.c != nil {
		c.c.Close()
		c.State = CONN_STATE_DISCONNECTED
	}
	return 0
}

func (c *TcpConn) StartRecv() {
	buf := make([]byte, 1024)
	for c.State == CONN_STATE_CONNECTED {
		n, err := c.c.Read(buf)
		if err != nil {
			c.Disconnect()
			break
		}

		rx <- buf[:n]
	}
}

func (c *TcpConn) StartSend() {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case buf := <-tx:
			_, err := c.c.Write(buf)
			if err != nil {
				c.Disconnect()
				break
			}
		}
	}
}

func (c *TcpConn) Read() int {
	n, err := c.c.Read(c.RxBuff)
	if err != nil {
		//read error, break the connection
		c.Disconnect()
		return -1
	}

	return n
}

func (c *TcpConn) ReadTo(buff []byte) int {
	n, err := c.c.Read(buff)
	if err != nil {
		//read error, break the connection
		c.Disconnect()
		return -1
	}
	return n
}

func (c *TcpConn) InstantWrite(buff []byte) int {
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
