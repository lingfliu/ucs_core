package conn

import (
	"net"
	"time"
)

type UdpConn struct {
	BaseConn

	c *net.UDPConn

	recv chan []byte
	send [][]byte
}

func (udp *UdpConn) Run() {
	conn, err := net.DialTimeout("udp", udp.RemoteAddr, time.Duration(udp.Timeout)*time.Millisecond)
	if err != nil {
		udp.State = CONN_STATE_DISCONNECTED
		udp.c = nil
		return
	}

	udp.c = conn.(*net.UDPConn)
	udp.State = CONN_STATE_CONNECTED
}

func (udp *UdpConn) Read(bs []byte) int {
	n, err := udp.c.Read(bs)
	if err != nil {
		return 0
	} else {
		return n
	}
}

func (udp *UdpConn) ReadToBuff() []byte {
	buff := make([]byte, 1024)
	n, err := udp.c.Read(buff)
	if err != nil {
		return nil
	} else {
		return buff[:n]
	}
}

func (udp *UdpConn) Write(bs []byte, n int) int {
	n, err := udp.c.Write(bs)
	if err != nil {
		return 0
	} else {
		return n
	}
}

func (udp *UdpConn) ScheduleWrite(bs []byte) {
	udp.send = append(udp.send, bs)
}

func (udp *UdpConn) Close() int {
	udp.State = CONN_STATE_DISCONNECTED
	udp.c.Close()
	return 0
}

func (udp *UdpConn) taskRead() {
	for udp.State == CONN_STATE_CONNECTED {
		buff := udp.ReadToBuff()
		udp.recv <- buff
	}
}

func (udp *UdpConn) taskWrite() {
	for udp.State == CONN_STATE_CONNECTED {
		bs := udp.send[0]
		n := udp.Write(bs, len(bs))
		if n > 0 {
		} else {
			udp.Close()
		}
	}
}
