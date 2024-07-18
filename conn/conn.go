package conn

import (
	"net"

	"github.com/lingfliu/ucs_core/utils"
)

const CONN_STATE_DISCONNECTED = 0
const CONN_STATE_CONNECTING = 1
const CONN_STATE_CONNECTED = 2

const CONN_CLASS_TCP string = "tcp"
const CONN_CLASS_UDP string = "udp"

// const CONN_CLASS_kcp string = "kcp" //TODO: implement kcp
const CONN_CLASS_QUIC string = "quic"
const CONN_CLASS_HTTP string = "http"

const CONN_CLASS_MQTT string = "mqtt"
const CONN_CLASS_RTSP string = "rtsp"

// ********************************************************
// conn bases
// ********************************************************
type BaseConn struct {
	State        int
	KeepAlive    bool  // by default true
	Timeout      int64 // connect timeout
	TimeoutRw    int64 // read write timeout
	LocalAddr    string
	RemoteAddr   string
	Port         int
	DisconnectAt int64
	ConnectedAt  int64
	Class        string //tcp, quic, http, mqtt, rtsp
}

type Conn interface {
	StartRecv(rx chan []byte)
	StartSend(tx chan []byte)
	Disconnect() int
	Connect() int
}

type TcpConn struct {
	BaseConn
	c *net.TCPConn
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
	addr := net.TCPAddr{
		IP:   net.ParseIP(c.RemoteAddr),
		Port: c.Port,
	}

	tcp, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		c.State = CONN_STATE_DISCONNECTED
		return -1
	}

	c.State = CONN_STATE_CONNECTED
	c.c = tcp

	return 0
}

func (c *TcpConn) Disconnect() int {
	if c.c != nil {
		c.c.Close()
		c.State = CONN_STATE_DISCONNECTED
	}

	return 0
}

func (c *TcpConn) StartRecv(rx chan []byte) {
	buf := make([]byte, 1024)
	for {
		n, err := c.c.Read(buf)
		if err != nil {
			c.Disconnect()
			break
		}

		rx <- buf[:n]
	}
}

func (c *TcpConn) StartSend(tx chan []byte) {
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

type TcpConnCli struct {
	Conn   *TcpConn
	rxBuff *utils.ByteRingBuffer
	txBuff *utils.ByteArrayRingBuffer
}

func (cli *TcpConnCli) Read() {
}

func (cli *TcpConnCli) ReadTo(buff []byte) {
}

func (cli *TcpConnCli) InstantWrite(buff []byte) {
}

func (cli *TcpConnCli) ScheduleWrite(buff []byte) {
	cli.txBuff.Curr()
}
