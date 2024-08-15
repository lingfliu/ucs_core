package conn

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

type UdpCli struct {
	RemoteAddr string
}

type UdpConn struct {
	BaseConn
	c *net.UDPConn

	cliSet     map[string]*UdpConn
	serverSide bool
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
		cliSet: make(map[string]*UdpConn),
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

	var udp net.Conn
	var err error
	var laddr string
	if len(strings.Split(c.RemoteAddr, ":")) > 1 {
		laddr = c.RemoteAddr
	} else {
		laddr = utils.IpPortJoin(c.RemoteAddr, c.Port)
	}

	udp, err = net.DialTimeout("udp", laddr, time.Duration(c.Timeout)*time.Millisecond)

	if err != nil {
		ulog.Log().I("udpconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port))
		c.State = CONN_STATE_DISCONNECTED
		return -1
	} else {

		//renew conn and rw control
		ulog.Log().I("udpconn", fmt.Sprintf("connected to %s:%d", c.RemoteAddr, c.Port))
		udp.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Nanosecond))
		c.serverSide = false //mark the conn to avoid exceptional close
		c.State = CONN_STATE_CONNECTED
		c.Io <- CONN_STATE_CONNECTED
		c.c = udp.(*net.UDPConn)
		c.RemoteAddr = udp.RemoteAddr().String()

		// stop previous tasks & create new sig
		c.cancelRw()

		c.sigRw, c.cancelRw = context.WithCancel(context.Background())
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
	delete(c.cliSet, c.RemoteAddr)
	return 0
}

/*
 * Listen on port
 * @param ch channel to send new conn
 * @return 0 if success, -1 if failed
 */
func (c *UdpConn) Listen(sigRun context.Context, ctxCfg context.Context, ch chan Conn) {
	laddr, err := net.ResolveUDPAddr("udp4", utils.IpPortJoin("127.0.0.1", c.Port))
	if err != nil {
		ulog.Log().E("udpconn", "address format err")
		c.Close()
		return
	}
	cc, err := net.ListenUDP("udp", laddr)
	if err != nil {
		ulog.Log().E("udpconn", "listen failed")
		c.Close()
		return
	}

	for c.State != CONN_STATE_CLOSED {
		select {
		case <-sigRun.Done():
			return
		default:
			bs := make([]byte, 1024)
			n, remoteAddr, err := cc.ReadFromUDP(bs)

			if err != nil {
				continue
			}

			if _, ok := c.cliSet[remoteAddr.String()]; ok {
				cli := c.cliSet[remoteAddr.String()]
				cli.Rx <- bs[:n]
				cli.lastRecvAt = utils.CurrentTime()
			} else {
				//new udp conn
				udp := &UdpConn{
					BaseConn: BaseConn{
						RemoteAddr: remoteAddr.String(),
						Tx:         make(chan []byte, 32),
						Rx:         make(chan []byte, 32),
						Io:         make(chan int),
					},
					c:          cc,
					serverSide: true,
				}

				udp.State = CONN_STATE_CONNECTED
				udp.lastConnectAt = utils.CurrentTime()
				udp.sigRun, udp.cancelRun = context.WithCancel(context.Background())
				udp.sigRw, udp.cancelRw = context.WithCancel(context.Background())
				// go udp._task_recv(udp.sigRw)
				go udp._task_send(udp.sigRw)

				c.cliSet[remoteAddr.String()] = udp
				ch <- udp
				udp.Rx <- bs[:n]

			}

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

	if c.c != nil && !c.serverSide {
		c.c.Close()
	}
	delete(c.cliSet, c.RemoteAddr)
	return 0
}

func (c *UdpConn) Write(bs []byte) int {
	if c.State != CONN_STATE_CONNECTED {
		return -1
	}

	var err error
	laddr, err := net.ResolveUDPAddr("udp", c.RemoteAddr)
	if err != nil {
		ulog.Log().E("udpconn", "address resolve error")
		c.Disconnect()
	}

	var n int
	if c.serverSide {
		n, err = c.c.WriteToUDP(bs, laddr)
	} else {
		c.c.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
		n, err = c.c.Write(bs)
	}
	if err != nil {
		ulog.Log().E("udpconn", "write failed"+err.Error())
		c.Disconnect()
		return -1
	}
	return n

}

func (c *UdpConn) _task_recv(sigRw context.Context) {
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
func (c *UdpConn) Read(bs []byte) int {
	c.c.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
	n, err := c.c.Read(bs)
	if err != nil {
		ulog.Log().E("udpconn", "read failed")
		c.Disconnect()
		return -1
	}
	return n
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
