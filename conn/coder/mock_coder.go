package coder

import (
	"github.com/lingfliu/ucs_core/utils"
)

type MockCoder struct {
	BaseCoder
}

/** A stub mock coder for dev and test
 *
 */
func NewMockCoder() MockCoder {
	c := MockCoder{
		BaseCoder: BaseCoder{
			proto:    "mock",
			Buff:     utils.NewByteRingBuffer(1024),
			codebook: nil,
		},
	}
	return c
}

func (c MockCoder) Encode(msg *Msg, bs []byte) int {
	return 1
}

func (c MockCoder) FastDecode(bs []byte) *Msg {
	msg := &Msg{
		Class:     0,
		Metas:     map[string]*MsgAttr{},
		Payloads:  map[string]*MsgAttr{},
		Timestamp: utils.CurrentTime(),
	}

	vals := make([]any, 1)
	vals[0] = string(bs)
	msg.Payloads["data"] = &MsgAttr{
		Name:     "data",
		ValClass: ATTR_CLASS_BYTE,
		Vals:     vals,
	}
	return msg
}

func (c MockCoder) Reset() {
	c.Buff.Flush()
}

func (c MockCoder) PushDecode(bs []byte, n int) *Msg {
	//TODO implement decoding from the buff
	return c.FastDecode(bs[:n])
}
