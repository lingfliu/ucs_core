package dd

import (
	"context"

	"github.com/lingfliu/ucs_core/conn"
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

/**
 * currently deprecated
 */
type MemDdSrv struct {
	Qos       int
	CliSet    map[string]*MemDdCli
	TopicSet  map[uint64]bool
	SubSet    map[uint64]map[string]MemDdCli
	State     int
	RxMsg     chan *DdZeroMsg
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
		case <-dd.RxMsg:
		// case ddMsg := <-dd.RxMsg:
		// if _, ok := dd.TopicSet[ddMsg.topic]; !ok {
		// 	continue
		// }

		// for _, cli := range dd.SubSet[ddMsg.topic] {
		// 	cli.Coder.TxMsg <- ddMsg
		// }

		case <-dd.sigRun.Done():
			return
		}
	}
}

type MemDdCli struct {
	ConnCli conn.ConnCli
	Coder   *DdZeroCoder
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
