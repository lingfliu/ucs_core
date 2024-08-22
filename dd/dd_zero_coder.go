package dd

import (
	"bytes"
	"context"

	"github.com/lingfliu/ucs_core/utils"
)

type DdZeroCoder struct {
	Buffer    *utils.ByteRingBuffer
	TxMsg     chan *DdZeroMsg
	RxMsg     chan *DdZeroMsg
	sigRun    context.Context
	cancelRun context.CancelFunc
	buff      []byte
}

func NewDdZeroCoder(buff_len int) *DdZeroCoder {
	sigRun, cancelRun := context.WithCancel(context.Background())
	coder := &DdZeroCoder{
		Buffer:    utils.NewByteRingBuffer(buff_len),
		TxMsg:     make(chan *DdZeroMsg, 1024),
		RxMsg:     make(chan *DdZeroMsg, 1024),
		sigRun:    sigRun,
		cancelRun: cancelRun,
		buff:      make([]byte, 1024),
	}
	return coder

}

func (coder *DdZeroCoder) Decode() *DdZeroMsg {
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
	// topic := utils.Byte2Int(coder.buff, 4, 4, true, true)
	msgLen := utils.Byte2Int(coder.buff, 8, 10, false, true)
	if coder.Buffer.Capacity < msgLen+10 {
		return nil
	} else {
		coder.Buffer.Pop(coder.buff, 10+msgLen)
		// msg := &DdMsg{
		// topic: uint64(topic),
		// data:  coder.buff[10 : 10+msgLen],
		// }

		return nil
	}
}

func (coder *DdZeroCoder) Encode(msg *DdZeroMsg) []byte {
	return nil
}

func (coder *DdZeroCoder) FastDecode(data []byte) *DdZeroMsg {
	return nil
}

func (coder *DdZeroCoder) Push(data []byte) {
	coder.Buffer.Push(data, len(data))
}

func (coder *DdZeroCoder) _task_decode() {
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
