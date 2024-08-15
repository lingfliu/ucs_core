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
	state := c.State
	if state == CONN_STATE_CONNECTED {
		return -2
	}
	if state == CONN_STATE_CONNECTING {
		return -2
	}

	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	udp, err := net.DialTimeout("udp", utils.IpPortJoin(c.RemoteAddr, c.Port), time.Duration(c.Timeout)*time.Millisecond)

	if err != nil {
		ulog.Log().I("udpconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED
		return -1
	} else {
		//renew conn and rw control
		ulog.Log().I("tcpconn", fmt.Sprintf("connected to %s:%d", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_CONNECTED
		c.Io <- CONN_STATE_CONNECTED
		c.c = udp.(*net.UDPConn)
		c.sigRw, c.cancelRw = context.WithCancel(context.Background())

		// stop previous tasks
		c.cancelRw()
		go c._task_recv(c.sigRw)
		go c._task_send(c.sigRw)
		return 0
	}
}

func (c *UdpConn) Disconnect() int {
	if c.State == CONN_STATE_DISCONNECTED || c.State == CONN_STATE_CLOSED {
		return -2
	}
	c.State = CONN_STATE_DISCONNECTED
	c.lastDisconnectAt = utils.CurrentTime()

	c.Io <- CONN_STATE_DISCONNECTED

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
func (c *UdpConn) Listen(sigRun context.Context, ctxCfg context.Context, ch chan Conn) {
	for c.State != CONN_STATE_CLOSED {
		select {
		case <-sigRun.Done():
			return
		default:
			laddr, err := net.ResolveUDPAddr("udp", utils.IpPortJoin("127.0.0.1", c.Port))
			if err != nil {
				ulog.Log().E("udpconn", "address format err")
				c.Close()
				return
			}

			cc, err := net.ListenUDP("udp", laddr)
			if err != nil {
				// ulog.Log().E("udpconn", "listen failed")
				continue
			}

			cfg := ctxCfg.Value(utils.CtxKeyCfg{}).(*ConnCfg)
			udp := NewUdpConn(cfg)
			ulog.Log().I("udpconn", "new udpconn, remoteaddr: ")
			// udp.RemoteAddr = cc.RemoteAddr().String()
			udp.c = cc
			udp.State = CONN_STATE_CONNECTED
			udp.lastConnectAt = utils.CurrentTime()
			udp.sigRun, udp.cancelRun = context.WithCancel(context.Background())
			udp.sigRw, udp.cancelRw = context.WithCancel(context.Background())
			go udp._task_recv(udp.sigRw)
			go udp._task_send(udp.sigRw)

			ch <- udp
		}
	}
}

func (c *UdpConn) Start(sigRun context.Context) chan int {
	go c._task_connect(sigRun)
	return c.Io
}

func (c *UdpConn) Close() int {
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
				c.Write(buff)
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
