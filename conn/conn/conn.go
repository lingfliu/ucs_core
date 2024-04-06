package conn

const CONN_STATE_DISCONNECTED = 1
const CONN_STATE_CONNECTED = 0

const CONN_CLASS_TCP string = "tcp"
const CONN_CLASS_UDP string = "udp"
const CONN_CLASS_kcp string = "kcp"
const CONN_CLASS_QUIC string = "quic"
const CONN_CLASS_HTTP string = "http"

// mqtt server
const CONN_CLASS_MQTT string = "mqtt"

// streaming servers
const CONN_CLASS_RTSP string = "rtsp"

// srv state
const SRV_STATE_ON = 0
const SRV_STATE_OFF = 1

// ********************************************************
// conn
// ********************************************************
type BaseConn struct {
	State        int
	KeepAlive    bool
	Timeout      int64 //connect timeout
	TimeoutRw    int64 //read write timeout
	LocalAddr    string
	RemoteAddr   string
	DisconnectAt int64
	ConnectedAt  int64
	Class        string //tcp, quic, http, mqtt, rtsp

	//callbacks
	OnRecv         func(bs []byte, n int)
	OnSent         func(bs []byte, n int)
	OnStateChanged func(state int)
}

type Conn interface {
	Run()

	ReadToBuff() ([]byte, int)
	Read(bs []byte) int
	Write(bs []byte) int
	ScheduledWrite(bs []byte)
	Disconnect() int
	Connect() int

	taskRead()
	taskWrite()
}
