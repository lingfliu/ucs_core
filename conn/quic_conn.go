package conn

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
	"github.com/quic-go/quic-go"
)

type QuicConn struct {
	BaseConn
	c      quic.Connection
	stream quic.Stream
}

func NewQuicConn(cfg *ConnCfg) *QuicConn {
	c := &QuicConn{
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

func (c *QuicConn) Connect() int {
	var err error
	var stream quic.Stream
	var qc quic.Connection

	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"ucs-quic"},
	}

	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	qc, err = quic.DialAddr(context.Background(), c.RemoteAddr, tlsCfg, nil)
	if err != nil {
		ulog.Log().I("quicconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED
		return -1
	}

	stream, err = qc.OpenStreamSync(context.Background())
	if err != nil {
		ulog.Log().I("quicconn", "stream opening failed")
		c.State = CONN_STATE_DISCONNECTED
		return -1
	}
	c.c = qc
	c.stream = stream
	c.State = CONN_STATE_CONNECTED

	go c._task_recv()
	go c._task_send()

	return 0
}

func (c *QuicConn) Disconnect() int {
	c.lastDisconnectAt = utils.CurrentTime()
	c.State = CONN_STATE_DISCONNECTED
	if c.c != nil {
		err := c.c.CloseWithError(0, "")
		if err != nil {
			ulog.Log().I("quicconn", "conn close error")
			return -1
		}
		err = c.stream.Close()
		if err != nil {
			ulog.Log().I("quicconn", "stream close error")
			return -1
		}
		return 0
	} else {
		return -1
	}
}

func (c *QuicConn) Listen(ch chan Conn) {
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: c.Port,
	}
	udpconn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		ulog.Log().E("quic", "listen failed, check port")
		c.Close()
		return
	}
	tr := &quic.Transport{Conn: udpconn}
	ln, err := tr.Listen(&tls.Config{InsecureSkipVerify: true}, &quic.Config{})

	for {
		cc, err := ln.Accept(context.Background())
		if err != nil {
			ulog.Log().E("tcpconn", "listen failed, check port")
			c.Close()
			break
		}
		qc := &QuicConn{
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
		ch <- qc //TODO: test the channel for new conn handling
	}
}

func (c *QuicConn) GetRxBuff() *utils.ByteRingBuffer {
	return c.RxBuff
}

func (c *QuicConn) StartRecv() {
	go c._task_recv()
}

func (c *QuicConn) _task_recv() {
	buff := make([]byte, 1024)
	for {
		c.stream.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
		n, err := c.stream.Read(buff)

		if err != nil {
			//TODO: handling disconnect
			c.Disconnect()
		}
		if n > 0 {
			c.RxBuff.Push(buff, n)
		}
	}
}

func (c *QuicConn) Read(bs []byte) int {
	c.stream.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
	n, err := c.stream.Read(bs)
	if err != nil {
		return 0
	} else {
		return n
	}
}

func (c *QuicConn) InstantWrite(bs []byte) int {
	err := c.stream.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
	if err != nil {
		return -1
	}

	n, err := c.stream.Write(bs)
	if err != nil {
		c.Disconnect()
		return -1
	} else {
		return n
	}
}

func (c *QuicConn) ScheduleWrite(bs []byte) {
	c.TxBuff.Push(bs)
}

func (c *QuicConn) _task_send() {
	for c.State == CONN_STATE_CONNECTED {
		buff := c.TxBuff.Pop()
		if buff == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		err := c.stream.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
		if err != nil {
			c.Disconnect()
		}

		_, err = c.stream.Write(buff)
		if err != nil {
			c.Disconnect()
		}
	}
}

func (c *QuicConn) _task_connect() {
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

func (c *QuicConn) Close() {
	c.State = CONN_STATE_CLOSE
	if c.c != nil {
		c.c.CloseWithError(0, "")
		c.stream.Close()
	}
}
