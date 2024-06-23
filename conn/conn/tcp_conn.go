package conn

import (
	"net"
	"time"

	"github.com/lingfliu/ucs_core/utils"
)

/**
 * TcpConn wrapper with improved read & write behavior
 * note: TcpConn only handle bytes read & write, connection state and decode are controlled by cli
 * method: Write, ScheduledWrite, Connect, Run, Disconnect
 */
type TcpConn struct {
	BaseConn

	c         *net.TCPConn
	recv_buff *utils.ByteArrayRingBuffer
	send_buff [][]byte
}

func NewTcpConn(localAddr string, remoteAddr string, keepAlive bool, timeout int64, timeoutRw int64) *TcpConn {
	if localAddr == "" {
		localAddr = "127.0.0.1"
	}

	c := &TcpConn{
		BaseConn: BaseConn{
			LocalAddr:  localAddr,
			RemoteAddr: remoteAddr,
			KeepAlive:  keepAlive,
			Timeout:    timeout,
			TimeoutRw:  timeoutRw,
			State:      CONN_STATE_DISCONNECTED,
			Class:      CONN_CLASS_TCP,
		},
		recv_buff: utils.NewByteArrayRingBuffer(10, 2048),
		send_buff: make([][]byte, 0),
	}
	return c
}

func (conn *TcpConn) taskRead() {

	// tick := time.NewTicker(100 * time.Microsecond)
	for conn.State == CONN_STATE_CONNECTED {

		// for range tick.C {
		// 	// buff, n = conn.Read()
		// 	// if n > 0 {
		// 	// 	if conn.OnRecv != nil {
		// 	// 		conn.OnRecv(buff, n)
		// 	// 	}
		// 	// }
		// }
	}
}

func (conn *TcpConn) taskWrite() {
	tick := time.NewTicker(100 * time.Microsecond)
	for conn.State == CONN_STATE_CONNECTED {
		for range tick.C {
			for len(conn.send_buff) > 0 {
				//TODO: replace send_buff with a FIFO queue
			}
		}
	}
}

func (conn *TcpConn) Read() ([]byte, int) {
	buff := conn.recv_buff.Next()
	//set ddl for every read
	err := conn.c.SetReadDeadline(time.Now().Add(time.Duration(conn.TimeoutRw) * time.Millisecond))
	if err != nil {
		return nil, 0
	}

	n, err := conn.c.Read(buff)
	if err != nil {
		return nil, 0
	} else {
		return buff, n
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

/*
 * called in io routines
 */
func (conn *TcpConn) Connect(state chan int) int {
	res := 0
	conn.State = CONN_STATE_CONNECTING
	c, err := net.DialTimeout("tcp", conn.RemoteAddr, time.Duration(conn.Timeout)*time.Millisecond)
	if err != nil {
		conn.State = CONN_STATE_DISCONNECTED
		conn.c = nil
		res = -1
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
		res = 0
	}

	return res
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
func (conn *TcpConn) Disconnect(state chan int) {
	conn.State = CONN_STATE_DISCONNECTED
	conn.DisconnectAt = utils.CurrentTime()

	err := conn.c.Close()
	if err != nil {
		// log.GetULogger().I("close failed " + err.Error())
		// simply delete the connection
		conn.c = nil
	}

	state <- CONN_STATE_CONNECTED
}

func (conn *TcpConn) taskState(state chan int) {
	for {
		select {
		case s := <-state:
			if s == CONN_STATE_DISCONNECTED {
				// reconnect after 1s
			}
		}
	}
}
