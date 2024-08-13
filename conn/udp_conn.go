package conn

import (
	"context"
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
	sigRun, cancelRun := context.WithCancel(context.Background())
	sigRw, cancelRw := context.WithCancel(context.Background())
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
			Rx:             make(chan []byte, 32),
			Tx:             make(chan []byte, 32),
			Io:             make(chan int),

			sigRun:    sigRun,
			cancelRun: cancelRun,
			sigRw:     sigRw,
			cancelRw:  cancelRw,
		},
		c: nil,
	}
	return c
}

func (c *UdpConn) Connect() int {
	if c.State == CONN_STATE_CONNECTED || c.State == CONN_STATE_CONNECTING {
		return -2
	}

	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	udp, err := net.DialTimeout("udp", c.RemoteAddr, time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		ulog.Log().I("udpconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED
		return -1
	} else {

		c.c = udp.(*net.UDPConn)

		c.State = CONN_STATE_CONNECTED
		c.Io <- CONN_STATE_CONNECTED
		c.c = udp.(*net.UDPConn)

		// stop previous tasks
		c.cancelRw()

		go c._task_recv(c.sigRun)
		go c._task_send(c.sigRun)
		return 0
	}
}

func (c *UdpConn) Disconnect() int {
	if c.State == CONN_STATE_DISCONNECTED || c.State == CONN_STATE_CLOSED {
		return -2
	}
	c.lastDisconnectAt = utils.CurrentTime()
	c.State = CONN_STATE_DISCONNECTED

	c.Io <- CONN_STATE_DISCONNECTED

	if c.c != nil {
		c.c.Close()
	}
	return 0
}

func (c *UdpConn) Listen(sigRun context.Context, ctxCfg context.Context, ch chan Conn) {
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: c.Port,
	}

	for c.State != CONN_STATE_CLOSED {
		select {
		case <-sigRun.Done():
			return
		default:
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

					lastConnectAt:    utils.CurrentTime(),
					lastDisconnectAt: 0,
					lastRecvAt:       utils.CurrentTime(),

					Rx: make(chan []byte, 32),
					Tx: make(chan []byte, 32),
					Io: make(chan int),
				},
				c: cc,
			}
			ch <- udp //TODO: test the channel for new conn handling
		}
	}
}

func (c *UdpConn) Read(bs []byte) (n int) {
	if c.State != CONN_STATE_CONNECTED {
		return -1
	}

	c.c.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))

	n, err := c.c.Read(bs)
	if err != nil {
		ulog.Log().E("udpconn", "read failed")
		c.Disconnect()
		return -1
	}
	return
}

func (c *UdpConn) Write(bs []byte) int {
	if c.State != CONN_STATE_CONNECTED {
		return -1
	}

	c.c.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
	n, err := c.c.Write(bs)
	if err != nil {
		ulog.Log().E("udpconn", "write failed")
		c.Disconnect()
		return -1
	}
	return n
}

func (c *UdpConn) Close() int {
	if c.State == CONN_STATE_CLOSED {
		return -2
	}

	c.State = CONN_STATE_CLOSED
	if c.c != nil {
		c.c.Close()
	}
	c.Io <- CONN_STATE_CLOSED
	return 0
}

func (c *UdpConn) GetRemoteAddr() string {
	return c.RemoteAddr
}

func (c *UdpConn) GetState() int {
	return c.State
}

// tasks
func (c *UdpConn) _task_recv(ctx context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case <-ctx.Done():
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

func (c *UdpConn) _task_send(ctx context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case buff := <-c.Tx:
			if len(buff) == 0 {
				continue
			} else {
				n := c.Write(buff)
				if n == -1 {
					ulog.Log().E("udpconn", "write failed, disconnect")
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *UdpConn) _task_connect(ctx context.Context) {
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
