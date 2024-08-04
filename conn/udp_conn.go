package conn

import (
	"fmt"
	"net"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

type UdpConn struct {
	BaseConn
	c *net.UDPConn
}

func NewUdpConn(cfg *ConnCfg) *UdpConn {
	c := &UdpConn{
		BaseConn: BaseConn{
			State:          CONN_STATE_DISCONNECTED,
			RemoteAddr:     cfg.RemoteAddr,
			Port:           cfg.Port,
			Class:          CONN_CLASS_UDP,
			KeepAlive:      cfg.KeepAlive,
			ReconnectAfter: cfg.ReconnectAfter,
			Timeout:        cfg.Timeout,
			TimeoutRw:      cfg.TimeoutRw,
			RxBuff:         utils.NewByteRingBuffer(1024),
			TxBuff:         utils.NewByteArrayRingBuffer(32, 1024),
		},
	}
	return c
}

func (c *UdpConn) Connect() int {
	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	udp, err := net.DialTimeout("udp", c.RemoteAddr, time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		ulog.Log().I("tcpconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED
		return -1
	}

	c.State = CONN_STATE_CONNECTED
	c.c = udp.(*net.UDPConn)

	go c._task_recv()
	go c._task_send()

	return 0
}

func (c *UdpConn) Disconnect() int {
	c.lastDisconnectAt = utils.CurrentTime()
	c.State = CONN_STATE_DISCONNECTED
	if c.c != nil {
		c.c.Close()
		//TODO: finish rx & tx
		return 0
	} else {
		return -1
	}
}

func (c *UdpConn) Listen(ch chan Conn) {
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: c.Port,
	}
	for {
		cc, err := net.ListenUDP("udp", &addr)
		if err != nil {
			ulog.Log().E("tcpconn", "listen failed, check port")
			c.Close()
			break
		}
		udp := &UdpConn{
			BaseConn: BaseConn{
				State:      CONN_STATE_CONNECTED,
				Class:      CONN_CLASS_TCP,
				KeepAlive:  c.KeepAlive,
				Timeout:    c.Timeout,
				TimeoutRw:  c.TimeoutRw,
				LocalAddr:  c.LocalAddr,
				RemoteAddr: cc.RemoteAddr().String(),
				Port:       c.Port,
				RxBuff:     utils.NewByteRingBuffer(1024),
				TxBuff:     utils.NewByteArrayRingBuffer(32, 1024),
			},
			c: cc,
		}
		ch <- udp //TODO: test the channel for new conn handling
	}
}

func (c *UdpConn) GetRxBuff() *utils.ByteRingBuffer {
	return c.RxBuff
}

func (c *UdpConn) StartRecv() {
	go c._task_recv()
}

func (c *UdpConn) _task_recv() {
	buff := make([]byte, 1024)
	for c.State == CONN_STATE_CONNECTED {
		c.c.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
		n, err := c.c.Read(buff)

		if err != nil {
			c.Disconnect()
		}
		if n > 0 {
			c.RxBuff.Push(buff, n)
		}
	}
}

func (c *UdpConn) Read(buff []byte) int {
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

func (c *UdpConn) InstantWrite(buff []byte) int {
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

func (c *UdpConn) ScheduleWrite(buff []byte) {
	c.TxBuff.Push(buff)
}

func (c *UdpConn) _task_send() {
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

func (c *UdpConn) _task_connect() {
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

func (c *UdpConn) Close() {
	c.State = CONN_STATE_CLOSE
	if c.c != nil {
		c.c.Close()
	}
}

func (c *UdpConn) GetRemoteAddr() string {
	return c.RemoteAddr
}
