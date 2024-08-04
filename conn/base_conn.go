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
	// CONN_CLASS_kcp string = "kcp" //TODO: implement kcp
)

// ********************************************************
// common attrs for conn
// ********************************************************
type BaseConn struct {
	State      int
	KeepAlive  bool  // by default true
	Timeout    int64 // connect timeout
	TimeoutRw  int64 // read write timeout
	LocalAddr  string
	RemoteAddr string
	Port       int
	Class      string //tcp, quic, http, mqtt

	ReconnectAfter   int64
	lastRecvAt       int64
	lastConnectAt    int64
	lastDisconnectAt int64

	// RxBuff *utils.ByteRingBuffer
	TxBuff *utils.ByteArrayRingBuffer

	//event
	OnClose        func()
	OnConnected    func()
	OnDisconnected func()

	//Signals
	SigCtl chan any
}

type Conn interface {
	Connect()
	Disconnect()
	Close()

	// RW
	ScheduleWrite([]byte)
	InstantWrite([]byte) int
	StartRecv() chan []byte

	//Srv
	Listen(ch chan Conn)
	//Establish()

	//Attr fetch
	GetRemoteAddr() string
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
