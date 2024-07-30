package conn

type TcpConnCli struct {
	c             *TcpConn
	KeepAlive     bool
	LastRecvAt    int64
	LastConnectAt int64
	LastMsgAt     int64
}
