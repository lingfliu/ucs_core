package coder

import (
	"bytes"
	"context"

	"github.com/lingfliu/ucs_core/utils"
)

/**
 * Zero coder for raw bytes redirecting,
 */
type ZeroCoder struct {
	Header []byte

	RxMsg chan *ZeroMsg

	//sigs
	sigRun    context.Context
	cancelRun context.CancelFunc
}

func NewZeroCoder() *ZeroCoder {
	ctx, cancel := context.WithCancel(context.Background())
	return &ZeroCoder{
		Header:    []byte{0xaa, 0xaa, 0xaa, 0xaa},
		RxMsg:     make(chan *ZeroMsg),
		sigRun:    ctx,
		cancelRun: cancel,
	}
}

func (coder *ZeroCoder) Encode(msg *ZeroMsg, bs []byte) int {
	payloadLen := len(msg.Data)
	copy(bs, coder.Header)
	utils.Int2Byte(msg.Class, bs, 4, 2, false, true)
	utils.Int2Byte(payloadLen, bs, 6, 2, false, true)
	return 8 + len(msg.Data)
}

func (coder *ZeroCoder) FastDecode(bs []byte) {
	if bytes.Equal(bs[:4], coder.Header) {
		class := utils.Byte2Int(bs, 4, 2, false, true)
		var payloadLen int
		var msg *ZeroMsg
		if class == 0 {
			//no payload
			payloadLen = 0

			msg = &ZeroMsg{
				Class: utils.Byte2Int(bs, 6, 2, false, true),
				Data:  nil,
			}
		} else {
			payloadLen = utils.Byte2Int(bs, 6, 2, false, true)
			msg = &ZeroMsg{
				Class: class,
				Data:  bs[8 : 8+payloadLen],
			}
		}
		coder.RxMsg <- msg
	}
}

func (coder *ZeroCoder) StartDecode(rxBytes chan []byte) {
	for {
		select {
		case <-coder.sigRun.Done():
			return
		case bs := <-rxBytes:
			coder.FastDecode(bs)
		}
	}
}

func (coder *ZeroCoder) StopDecode() {
	coder.cancelRun()
}

func (coder *ZeroCoder) CreatePingpongMsg() *ZeroMsg {
	return &ZeroMsg{
		Class: 0,
		Data:  []byte{},
	}
}
