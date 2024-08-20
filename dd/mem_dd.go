package dd

import (
	"bytes"
	"context"

	"github.com/lingfliu/ucs_core/conn"
	"github.com/lingfliu/ucs_core/utils"
)

const (
	DD_STATE_ON   = 0
	DD_STATTE_OFF = 1
)

/**
 * int hash function for topic
 */
func topic2int(topic string) uint64 {
	return 0
}

type DdMsg struct {
	topic uint64
	data  []byte
}

type DdMsgCoder struct {
	Buffer    *utils.ByteRingBuffer
	TxMsg     chan *DdMsg
	RxMsg     chan *DdMsg
	sigRun    context.Context
	cancelRun context.CancelFunc
	buff      []byte
}

func NewDdMsgCoder(buff_len int) *DdMsgCoder {
	sigRun, cancelRun := context.WithCancel(context.Background())
	coder := &DdMsgCoder{
		Buffer:    utils.NewByteRingBuffer(buff_len),
		TxMsg:     make(chan *DdMsg, 1024),
		RxMsg:     make(chan *DdMsg, 1024),
		sigRun:    sigRun,
		cancelRun: cancelRun,
		buff:      make([]byte, 1024),
	}
	return coder

}

func (coder *DdMsgCoder) Decode() *DdMsg {
	//zero coder
	if coder.Buffer.Capacity < 6 {
		return nil
	}

	coder.Buffer.Peek(coder.buff, 4)
	for !bytes.Equal(coder.buff[:4], []byte{0xAA, 0xAA, 0xAA, 0xAA}) {
		coder.Buffer.Pop(coder.buff, 1)
		if coder.Buffer.Capacity < 6 {
			return nil
		}
	}

	coder.Buffer.Peek(coder.buff, 8)
	topic := utils.Byte2Int(coder.buff, 4, 4, true, true)
	msgLen := utils.Byte2Int(coder.buff, 8, 10, false, true)
	if coder.Buffer.Capacity < msgLen+10 {
		return nil
	} else {
		coder.Buffer.Pop(coder.buff, 10+msgLen)
		msg := &DdMsg{
			topic: uint64(topic),
			data:  coder.buff[10 : 10+msgLen],
		}

		return msg
	}
}

func (coder *DdMsgCoder) Encode(msg *DdMsg) []byte {
	return nil
}

func (coder *DdMsgCoder) FastDecode(data []byte) *DdMsg {
	return nil
}

func (coder *DdMsgCoder) Push(data []byte) {
	coder.Buffer.Push(data, len(data))
}

func (coder *DdMsgCoder) _task_decode() {
	for {
		select {
		case <-coder.sigRun.Done():
			return
		default:
			msg := coder.Decode()
			if msg != nil {
				coder.RxMsg <- msg
			}
		}
	}
}

type MemDdSrv struct {
	Qos       int
	CliSet    map[string]*MemDdCli
	TopicSet  map[uint64]bool
	SubSet    map[uint64]map[string]MemDdCli
	State     int
	RxMsg     chan *DdMsg
	sigRun    context.Context
	cancelRun context.CancelFunc
}

func (dd *MemDdSrv) Start() {
}

func (dd *MemDdSrv) Stop() {
}

func (dd *MemDdSrv) CliSubscribe(topic string, data []byte) {
	topicTag := topic2int(topic)
	if _, ok := dd.TopicSet[topicTag]; !ok {
		dd.TopicSet[topicTag] = true
	}
}

func (dd *MemDdSrv) task_recv() {
	for dd.State == DD_STATE_ON {
		select {
		case ddMsg := <-dd.RxMsg:
			if _, ok := dd.TopicSet[ddMsg.topic]; !ok {
				continue
			}

			for _, cli := range dd.SubSet[ddMsg.topic] {
				cli.Coder.TxMsg <- ddMsg
			}

		case <-dd.sigRun.Done():
			return
		}
	}
}

type MemDdCli struct {
	ConnCli conn.ConnCli
	Coder   *DdMsgCoder
	State   int

	//connection
	Host  string
	Prop  int
	Token string

	SubTopicSet map[int]bool

	sigRun context.Context
	context.CancelCauseFunc
}

func (cli *MemDdCli) Connect() int {
	ret := cli.ConnCli.Conn.Connect()
	if ret < 0 {
		return -1
	} else {
		return 0
	}
}

func (cli *MemDdCli) Disconnect() {
}

func (cli *MemDdCli) Subscribe(topic string) {
}

func (cli *MemDdCli) Start() {
	go cli._task_connect()
}

func (cli *MemDdCli) Close() {
}

func (cli *MemDdCli) _task_connect() {
	for cli.State == DD_STATE_ON {
		select {
		case <-cli.sigRun.Done():
			return
		}
	}
}
