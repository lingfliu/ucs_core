package conn

import (
	"net"
	"time"

	// "github.com/lingfliu/ucs_core/cfg"
	"github.com/lingfliu/ucs_core/utils"
)

/**
 * TcpConn wrapper with improved read & write behavior
 * note: TcpConn only handle bytes read & write, connection state and decode are controlled by the tcp_conn_cli
 * method: ReadToBuff, Read, Write, ScheduledWrite, Connect, Run, Disconnect
 */
type TcpConn struct {
	BaseConn

	c         *net.TCPConn
	recv_buff *utils.ByteArrayRingBuffer
	send_buff [][]byte
}

func NewTcpConn(localAddr string, remoteAddr string, keepAlive bool, timeout int64, timeoutRw int64) *TcpConn {
	var c = &TcpConn{
		BaseConn: BaseConn{
			LocalAddr:  localAddr,
			RemoteAddr: remoteAddr,
			KeepAlive:  keepAlive,
			Timeout:    timeout,
			TimeoutRw:  timeoutRw,
			State:      CONN_STATE_DISCONNECTED,
			Class:      CONN_CLASS_TCP,
		},
		recv_buff: utils.NewByteArrayRingBuffer(32, 2048),
		send_buff: make([][]byte, 0),
	}
	return c
}

func (conn *TcpConn) taskRead() {
	var n int
	var buff []byte
	tick := time.NewTicker(100 * time.Microsecond)
	for conn.State == CONN_STATE_CONNECTED {
		for range tick.C {
			buff, n = conn.ReadToBuff()
			if n > 0 {
				if conn.OnRecv != nil {
					conn.OnRecv(buff, n)
				}
			}
		}
	}
}

func (conn *TcpConn) taskWrite() {
	var n int
	tick := time.NewTicker(100 * time.Microsecond)
	for conn.State == CONN_STATE_CONNECTED {
		for range tick.C {
			for len(conn.send_buff) > 0 {
				//TODO: replace send_buff with a FIFO queue
				n = conn.Write(conn.send_buff[0])
				if n > 0 {
					if conn.OnSent != nil {
						conn.OnSent(conn.send_buff[0], n)
					}
					conn.send_buff = conn.send_buff[1:]
				}
			}
		}
	}
}

func (conn *TcpConn) ReadToBuff() ([]byte, int) {
	buff := conn.recv_buff.Next()
	conn.c.SetReadDeadline(time.Now().Add(time.Duration(conn.TimeoutRw) * time.Millisecond))
	n, _ := conn.c.Read(buff)
	if n > 0 {
		return buff, n
	} else {
		return nil, 0
	}
}

func (conn *TcpConn) Read(bs []byte) int {
	err := conn.c.SetReadDeadline(time.Now().Add(time.Duration(conn.TimeoutRw) * time.Millisecond))
	if err != nil {
		return 0
	}
	n, _ := conn.c.Read(bs)
	if (n > 0) && (conn.OnRecv != nil) {
		return n
	} else {
		return 0
	}
}

func (conn *TcpConn) Write(bs []byte) int {

	err := conn.c.SetWriteDeadline(time.Now().Add(time.Duration(conn.TimeoutRw) * time.Millisecond))
	if err != nil {
		return 0
	}
	n, err := conn.c.Write(bs)
	if err != nil {
		// log.GetULogger().I("write failed " + err.Error())
		return -1
	} else {
		return n
	}
}

func (conn *TcpConn) ScheduledWrite(bs []byte) {
	conn.send_buff = append(conn.send_buff, bs)
}

func (conn *TcpConn) Connect() int {
	c, err := net.DialTimeout("tcp", conn.RemoteAddr, time.Duration(conn.Timeout)*time.Millisecond)
	if err != nil {
		// cfg.GetULogger().E("connect error: " + err.Error())
		conn.State = CONN_STATE_DISCONNECTED
		conn.c = nil

		if conn.OnStateChanged != nil {
			conn.OnStateChanged(CONN_STATE_DISCONNECTED)
		}

		return -1
	} else {
		conn.c = c.(*net.TCPConn)

		//set tcp as nodelay
		conn.c.SetNoDelay(true)
		conn.c.SetLinger(0)
		conn.c.SetWriteBuffer(0)
		conn.c.SetReadBuffer(0)

		//initialize tcp state
		conn.State = CONN_STATE_CONNECTED
		conn.ConnectedAt = utils.CurrentTime()

		if conn.OnStateChanged != nil {
			conn.OnStateChanged(CONN_STATE_CONNECTED)
		}
		return 0
	}
}

/**
 * Start periodic tasks for read & write, only effective for keepAlive
 */
func (conn *TcpConn) Run() {
	if conn.KeepAlive {
		go conn.taskRead()
		go conn.taskWrite()
	}
}

/**
 * Disconnect whenever the op is successful or not
 */
func (conn *TcpConn) Disconnect() int {
	conn.State = CONN_STATE_DISCONNECTED
	conn.DisconnectAt = utils.CurrentTime()

	if conn.OnStateChanged != nil {
		conn.OnStateChanged(CONN_STATE_DISCONNECTED)
	}

	err := conn.c.Close()
	if err != nil {
		return -1
	} else {
		return 0
	}
}
