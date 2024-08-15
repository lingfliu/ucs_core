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
	sigRun, cancelRun := context.WithCancel(context.Background())
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
			Rx:             make(chan []byte, 32),
			Tx:             make(chan []byte, 32),
			Io:             make(chan int),
			sigRun:         sigRun,
			cancelRun:      cancelRun,
		},
		c: nil,
	}
	return c
}

func (c *QuicConn) Connect() int {
	if c.State != CONN_STATE_DISCONNECTED {
		return -2
	}

	var err error
	var stream quic.Stream
	var qc quic.Connection

	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"ucs-quic"},
	}

	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	qc, err = quic.DialAddr(context.Background(), utils.IpPortJoin(c.RemoteAddr, c.Port), tlsCfg, nil)
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
	c.Io <- CONN_STATE_CONNECTED

	go c._task_recv(c.sigRun)
	go c._task_send(c.sigRun)

	return 0
}

func (c *QuicConn) Disconnect() int {
	if c.State == CONN_STATE_CLOSED || c.State == CONN_STATE_DISCONNECTED {
		return -2
	}

	c.lastDisconnectAt = utils.CurrentTime()
	c.State = CONN_STATE_DISCONNECTED

	close(c.Rx)
	c.Io <- CONN_STATE_DISCONNECTED

	err := c.c.CloseWithError(0, "")
	if err != nil {
		ulog.Log().I("quicconn", "conn close error")
	}
	err = c.stream.Close()
	if err != nil {
		ulog.Log().I("quicconn", "stream close error")
	}
	return 0
}

func (c *QuicConn) Listen(sigRun context.Context, ctxCfg context.Context, ch chan Conn) {
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: c.Port,
	}
	udpconn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		ulog.Log().E("quicconn", "listen failed, check port")
		c.Close()
		return
	}
	tr := &quic.Transport{Conn: udpconn}
	ln, err := tr.Listen(&tls.Config{InsecureSkipVerify: true}, &quic.Config{})
	if err != nil {
		ulog.Log().E("quicconn", "listen failed, check config")
	}

	for c.State != CONN_STATE_CLOSED {
		select {
		case <-sigRun.Done():
			return
		default:
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
					Rx:         make(chan []byte, 32),
					Tx:         make(chan []byte, 32),
					Io:         make(chan int),
				},
				c: cc,
			}

			ch <- qc //TODO: test the channel for new conn handling
		}
	}
}

func (c *QuicConn) _task_recv(ctx context.Context) {
	buff := make([]byte, 1024)
	for c.State == CONN_STATE_CONNECTED {
		select {
		case run := <-ctx.Value("conn_trl").(chan bool):
			if !run {
				return
			}
		default:
			c.stream.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
			n, err := c.stream.Read(buff)

			if err != nil {
				//TODO: handling disconnect
				c.Disconnect()
				break
			}
			if n > 0 {
				c.Rx <- buff[:n]
			}
		}
	}
}

func (c *QuicConn) _task_send(ctx context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case buff := <-c.Tx:
			if len(buff) == 0 {
				// time.Sleep(100 * time.Millisecond)
				continue
			}
			err := c.stream.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
			if err != nil {
				c.Disconnect()
				return
			}
			_, err = c.stream.Write(buff)
			if err != nil {
				c.Disconnect()
				return
			}
		case run := <-ctx.Value("conn_ctrl").(chan bool):
			if !run {
				return
			}
		}
	}
}

func (c *QuicConn) _task_connect() {
	tic := time.NewTicker(time.Second * 1)
	for c.State != CONN_STATE_CLOSED {
		select {
		case <-tic.C:
			if c.State == CONN_STATE_DISCONNECTED && utils.CurrentTime()-c.lastDisconnectAt > c.ReconnectAfter {
				c.Connect()
			}
		}
	}
}

func (c *QuicConn) Close() int {
	if c.State == CONN_STATE_CLOSED {
		return -2
	}
	c.State = CONN_STATE_CLOSED
	if c.c != nil {
		c.c.CloseWithError(0, "")
		c.stream.Close()
	}
	c.Io <- CONN_STATE_CLOSED
	return 0
}

func (c *QuicConn) GetRemoteAddr() string {
	return c.RemoteAddr
}
