package conn

import (
	"time"

	"github.com/lingfliu/ucs_core/conn/coder"
)

const (
	CONN_CLI_STATE_CONNECTED    = 0
	CONN_CLI_STATE_DISCONNECTED = 1
	CONN_CLI_STATE_CONNECTING   = 2

	CONN_FLOW_STATE_WAIT    = 0
	CONN_FLOW_STATE_ACK     = 1
	CONN_FLOW_STATE_TIMEOUT = 2
)

type ConnFlow struct {
	ReqMsg     *coder.Msg
	AckTimeout int64
	State      int
	Ts         int64 // timestamp of req
}

type ConnCli struct {
	ConnFlowSet    map[int]*ConnFlow
	FlowTimeoutCnt int
	Coder          coder.Coder
	CoderName      string
	Conn           Conn
	ConnClass      string
	State          int
	TxMsgQueue     chan *coder.UMsg
}

func (cli *ConnCli) Connect() {
	res := cli.Conn.Connect()
	if res != 0 {

	} else {
		go cli.taskRead()
		go cli.taskWrite()
	}
}

func (cli *ConnCli) Disconnect() {

}

// func (cli *ConnCli) Submit(msg *coder.Msg) {
// 	cli.TxMsgQueue <- msg
// }

func (cli *ConnCli) Shutdown() {

}

func (cli *ConnCli) taskRead() {
	for {
		bs := cli.Conn.ReadToBuff()
		msg := cli.Coder.PushDecode(bs, len(bs))
		if msg != nil {
			//TODO: pass to channel
			_, ok := cli.ConnFlowSet[msg.Class]
			if ok {
				cli.ConnFlowSet[msg.Class].State = CONN_FLOW_STATE_WAIT
			}
		}
	}
}

func (cli *ConnCli) taskWrite() {
	for msg := range cli.TxMsgQueue {
		bs := make([]byte, 1024)
		cli.Coder.Encode(msg, bs)
		cli.Conn.Write(bs)
	}
}

func (cli *ConnCli) Submit(msg *coder.UMsg, ackTimeout int64) {
	//TODO: submit msg

	go cli.taskAckTimeout(msg, ackTimeout)
}

func (cli *ConnCli) taskAckTimeout(msg *coder.UMsg, ackTimeout int64) {
	//wait for timeout
	tic := time.NewTicker(time.Duration(ackTimeout) * (time.Millisecond))
	_ = <-tic.C
	if cli.ConnFlowSet[msg.Class].State == CONN_FLOW_STATE_ACK {
		//TODO: handling timeout
		cli.FlowTimeoutCnt += 1

	}
}

func (cli *ConnCli) taskLinkCheck() {
	tic := time.NewTicker(2 * time.Second)
	for range tic.C {
		//TODO: replace 3 with configurable value
		if cli.FlowTimeoutCnt > 3 {
			cli.Reconnect(2000)
		}
	}
}

func (cli *ConnCli) Reconnect(timeout int64) {
	cli.Disconnect()
	tic := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	_ = <-tic.C
	cli.Connect()
}
