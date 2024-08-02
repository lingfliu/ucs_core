package conn

import "github.com/lingfliu/ucs_core/utils"

const (
	CONN_STATE_DISCONNECTED = 0
	CONN_STATE_CONNECTING   = 1
	CONN_STATE_CONNECTED    = 2
	CONN_STATE_CLOSE        = 3

	CONN_CLASS_TCP  string = "tcp"
	CONN_CLASS_UDP  string = "udp"
	CONN_CLASS_QUIC string = "quic"
	CONN_CLASS_HTTP string = "http"
	CONN_CLASS_MQTT string = "mqtt"
	CONN_CLASS_RTSP string = "rtsp"
)

// const CONN_CLASS_kcp string = "kcp" //TODO: implement kcp

// ********************************************************
// conn bases
// ********************************************************
type BaseConn struct {
	State      int
	KeepAlive  bool  // by default true
	Timeout    int64 // connect timeout
	TimeoutRw  int64 // read write timeout
	LocalAddr  string
	RemoteAddr string
	Port       int
	Class      string //tcp, quic, http, mqtt, rtsp

	ReconnectAfter   int64
	lastRecvAt       int64
	lastConnectAt    int64
	lastDisconnectAt int64

	RxBuff *utils.ByteRingBuffer
	TxBuff *utils.ByteArrayRingBuffer

	//event
	OnConnected    func()
	OnDisconnected func()
}

type Conn interface {
	Connect() int
	Disconnect() int

	ScheduleWrite([]byte)
	InstantWrite([]byte) int

	Listen(ch chan Conn)
	Close()

	//Ops
	GetRxBuff() *utils.ByteRingBuffer
	StartRecv()
}

type ConnCfg struct {
	RemoteAddr     string
	LocalAddr      string
	Port           int
	Class          string
	KeepAlive      bool
	ReconnectAfter int64
	Timeout        int64
	TimeoutRw      int64
}
