package conn

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
	"github.com/quic-go/quic-go"
)

type QuicConn struct {
	BaseConn
	c      quic.Connection
	tlsCfg *tls.Config
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
		tlsCfg: generateTLSConfig(),
		c:      nil,
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
		NextProtos:         []string{"ucs"},
	}

	c.State = CONN_STATE_CONNECTING
	c.lastConnectAt = utils.CurrentTime()

	ctx := context.Background() //.WithCancel(context.Background())
	// defer cancel()
	qc, err = quic.DialAddr(ctx, utils.IpPortJoin(c.RemoteAddr, c.Port), tlsCfg, nil)
	if err != nil {
		ulog.Log().I("quicconn", fmt.Sprintf("connect to %s:%d failed", c.RemoteAddr, c.Port)+" "+err.Error())
		c.State = CONN_STATE_DISCONNECTED
		return -1
	}

	stream, err = qc.OpenStreamSync(context.Background())
	if err != nil {
		ulog.Log().I("quicconn", "stream opening failed err: "+err.Error())
		c.State = CONN_STATE_DISCONNECTED
		return -1
	}

	c.c = qc
	c.stream = stream

	c.State = CONN_STATE_CONNECTED
	c.Io <- CONN_STATE_CONNECTED

	sigRw, cancelRw := context.WithCancel(context.Background())
	c.sigRw = sigRw
	c.cancelRw = cancelRw

	go c._task_recv(c.sigRw)
	go c._task_send(c.sigRw)

	return 0
}

func (c *QuicConn) Start(sigRun context.Context) chan int {
	go c._task_connect(sigRun)
	return c.Io
}

func (c *QuicConn) Disconnect() int {
	if c.State == CONN_STATE_CLOSED || c.State == CONN_STATE_DISCONNECTED {
		return -2
	}

	c.lastDisconnectAt = utils.CurrentTime()
	c.State = CONN_STATE_DISCONNECTED

	c.Io <- CONN_STATE_DISCONNECTED
	c.cancelRw()

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

	// udpConn, err := net.ListenUDP("udp", &net.UDPAddr{Port: c.Port})
	// if err != nil {
	// 	ulog.Log().E("quicconn", "listen failed, check port")
	// 	c.Close()
	// 	return
	// }
	// tr := &quic.Transport{Conn: udpConn}

	ln, err := quic.ListenAddr(utils.IpPortJoin("0.0.0.0", c.Port), generateTLSConfig(), nil)
	// ln, err := tr.Listen(generateTLSConfig(), nil)
	if err != nil {
		ulog.Log().E("quicconn", "listen failed, check config")
		c.Close()
	}

	for c.State != CONN_STATE_CLOSED {
		select {
		case <-sigRun.Done():
			return
		default:
			cc, err := ln.Accept(context.Background())
			if err != nil {
				ulog.Log().E("quicconn", "listen failed, check port")
				continue
			}

			stream, err := cc.AcceptStream(context.Background())
			// stream, err := cc.AcceptStream(context.Background())
			io.Copy(stream, stream)

			// if err != nil {
			// 	continue
			// }
			// cfg := ctxCfg.Value(utils.CtxKeyCfg{}).(*ConnCfg)
			// qc := NewQuicConn(cfg)
			// qc.c = cc
			// qc.stream = stream
			// qc.RemoteAddr = cc.RemoteAddr().String()
			// qc.State = CONN_STATE_CONNECTED
			// qc.sigRun, qc.cancelRun = context.WithCancel(context.Background())
			// qc.sigRw, qc.cancelRw = context.WithCancel(context.Background())
			// go qc._task_recv(qc.sigRw)
			// go qc._task_send(qc.sigRw)

			// ch <- qc
		}
	}
}

func (c *QuicConn) _task_recv(sigRw context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case <-sigRw.Done():
			return
		default:
			buff := make([]byte, 1024)
			c.stream.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Nanosecond))
			n, err := c.stream.Read(buff)
			ulog.Log().I("quicconn", "recv len: "+strconv.Itoa(n))
			if err != nil {
				if err == io.EOF {
					c.Disconnect()
					return
				}
				//do not handle timeout here
				continue
			}
			if n > 0 {
				c.Rx <- buff[:n]
			}
		}
	}
}

func (c *QuicConn) _task_send(sigRw context.Context) {
	for c.State == CONN_STATE_CONNECTED {
		select {
		case buff := <-c.Tx:
			if len(buff) == 0 {
				continue
			}
			err := c.stream.SetWriteDeadline(time.Now().Add(time.Duration(c.TimeoutRw) * time.Millisecond))
			if err != nil {
				c.Disconnect()
				return
			}
			n, err := c.stream.Write(buff)
			if err != nil {
				ulog.Log().E("quicconn", "write failed")
				c.Disconnect()
				return
			}

			ulog.Log().E("quicconn", "writen to stream len: "+strconv.Itoa(n))
		case <-sigRw.Done():
			return
		}
	}
}

func (c *QuicConn) _task_connect(sigRun context.Context) {
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

func (c *QuicConn) Close() int {
	if c.State == CONN_STATE_CLOSED {
		return -2
	}
	c.State = CONN_STATE_CLOSED
	if c.c != nil {
		c.c.CloseWithError(0, "")
	}
	if c.stream != nil {
		c.stream.Close()
	}
	c.Io <- CONN_STATE_CLOSED
	return 0
}

func (c *QuicConn) GetRemoteAddr() string {
	return c.RemoteAddr
}

// copied from the quic-go example
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := func() []byte {
		var buf bytes.Buffer
		if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
			return nil
		}
		return buf.Bytes()
	}()

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"ucs"},
	}
}
