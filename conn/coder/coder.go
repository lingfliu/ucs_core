package coder

import (
	"github.com/lingfliu/ucs_core/utils"
)

// some reserved attr
const ATTR_CLASS_INT = 0
const ATTR_CLASS_FLOAT = 1
const ATTR_CLASS_BOOL = 2
const ATTR_CLASS_BYTE = 3

/**
 * define coder base
 * proto: protocol name, by default is "general_declare" which follows the reference declarative protocol
 */
type BaseCoder struct {
	proto    string
	Buff     *utils.ByteRingBuffer
	codebook *Codebook
}

type Coder interface {
	Reset()
	Encode(msg *Msg, bs []byte) int   //encode msg to byte, return length of the msg
	PushDecode(bs []byte, n int) *Msg //push bs into the coder and try to decode from the ringbuffer return nil if decaode fails
	FastDecode(bs []byte) *Msg        //fast decode without passing through ring buffer
}
