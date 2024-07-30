package conn

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

type ConnCli interface {
	Disconnect() int
	Connect() int

	StartRecv(rx chan []byte)
	StartSend(tx chan []byte)
}

type ConnSrv interface {
	Listen() int
	Disconnect() int

	StartRecv(rx chan []byte)
	StartSend(tx chan []byte)

	Cleanup()
}
