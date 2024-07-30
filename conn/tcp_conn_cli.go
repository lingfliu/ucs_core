package conn

import (
	"github.com/lingfliu/ucs_core/utils"
)

type TcpConnCli struct {
	c              *TcpConn
	KeepAlive      bool
	ReconnectAfter int64

	lastRecvAt       int64
	lastConnectAt    int64
	lastDisconnectAt int64
	lastMsgAt        int64
}

func NewTcpConnCli(remoteAddr string, port int, keepAlive bool, ReconnectAfter int64) *TcpConnCli {
	c := &TcpConnCli{
		c:                NewTcpConn(remoteAddr, port),
		lastRecvAt:       0,
		lastConnectAt:    0,
		lastDisconnectAt: 0,
		lastMsgAt:        0,
		KeepAlive:        keepAlive,
		ReconnectAfter:   ReconnectAfter,
	}
	return c
}

func (cli *TcpConnCli) Connect() int {
	cli.lastConnectAt = utils.CurrentTime()
	return cli.c.Connect()
}

func (cli *TcpConnCli) Disconnect() int {
	cli.lastDisconnectAt = utils.CurrentTime()
	return cli.c.Disconnect()
}

func (cli *TcpConnCli) ScheduleWrite(data []byte) {
	cli.c.ScheduleWrite(data)
}

func (cli *TcpConnCli) InstantWrite(data []byte) {
	cli.c.InstantWrite(data)
}

func (cli *TcpConnCli) GetRxBuff() *utils.ByteRingBuffer {
	return cli.c.RxBuff
}
