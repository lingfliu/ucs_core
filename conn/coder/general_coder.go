package coder

import (
	"bytes"

	"github.com/lingfliu/ucs_core/utils"
)

/**
 * declarative protocol decoder
 */
type GeneralCoder struct {
	BaseCoder

	tmpBuff  []byte
	codebook *Codebook
}

func NewGeneralCoder(codebook *Codebook) *GeneralCoder {
	c := &GeneralCoder{
		BaseCoder: BaseCoder{
			proto:    "deca",
			Buff:     utils.NewByteRingBuffer(1024),
			codebook: codebook,
		},
		codebook: codebook,
		tmpBuff:  make([]byte, 1024),
	}
	return c
}

func (c *GeneralCoder) FastDecode(bs []byte) *Msg {
	return nil
}

/**
 *
 */
//TODO decode zero-payload message
func (c *GeneralCoder) PushDecode(bs []byte, n int) *Msg {
	cb := c.codebook
	header := cb.Header

	for c.Buff.Capacity > len(header) {
		c.Buff.Peek(c.tmpBuff, len(header))
		if bytes.Equal(c.tmpBuff[:len(header)], header) {
			c.Buff.Drop(1)
			continue
		}

		cco, meta := c.codebook.GetCodeClassOffset()
		if cco+meta.Dimen > c.Buff.Capacity {
			// no enough bytes, wait
			return nil
		}

		class := c.decodeCodeClass(cco, meta)

		if !c.codebook.HasClass(class) {
			//invalid class, drop bytes
			c.Buff.Drop(1)
			return nil
		}

		if c.codebook.Varlen(class) {
			payloadSize, metaSize := c.estimatePayloadSize(class)
			if payloadSize < 0 {
				// no enough bytes to get all payload length, wait
				return nil
			} else {
				// found all payload size
				if payloadSize+metaSize+len(header) > c.Buff.Capacity {
					// no enough bytes, wait
					return nil
				} else {
					//TODO validate erc

					attrMetas := map[string]*MsgAttr{}
					//first decode meta
					for _, meta := range c.codebook.GetMetas(class) {
						attrMeta := c.decodeMeta(meta)
						attrMetas[meta.Name] = attrMeta
					}

					//then decode payload (with variable length)
					payloadMetas := c.codebook.GetPayloads(class)
					attrPayloads := map[string]*MsgAttr{}
					for _, meta := range payloadMetas {
						attrPayload := c.decodeMeta(meta)
						attrPayloads[meta.Name] = attrPayload
					}

					msg := &Msg{
						Class:    class,
						Metas:    attrMetas,
						Payloads: attrPayloads,
					}

					return msg
				}
			}
		} else {
			payloadSize, metaSize := c.getSize(class)
			if payloadSize > 0 && payloadSize+metaSize+len(header) > c.Buff.Capacity {
				// no enough bytes, wait
				return nil
			} else {
				//TODO validate erc

				attrMetas := map[string]*MsgAttr{}
				//first decode meta
				for _, meta := range c.codebook.GetMetas(class) {
					attrMetas[meta.Name] = c.decodeMeta(meta)
				}

				//then decode payload (with variable length)
				payloadMetas := c.codebook.GetPayloads(class)
				attrPayloads := map[string]*MsgAttr{}
				for _, meta := range payloadMetas {
					attrPayload := c.decodeMeta(meta)
					attrPayloads[meta.Name] = attrPayload
				}

				msg := &Msg{
					Class:    class,
					Metas:    attrMetas,
					Payloads: attrPayloads,
				}
				return msg
			}
		}

	}
	return nil
}

// TODO !!! implement the following functions !!!
func (c *GeneralCoder) Reset() {
	c.Buff.Flush()
}

func (c *GeneralCoder) Encode(msg *Msg, bs []byte) int {
	return 0
}

func (c *GeneralCoder) decodeCodeClass(offset int, meta *CodeMeta) int {
	return 0
}

/**
 * estimate payload size, if buff is not enough, return -1
 */
func (c *GeneralCoder) estimatePayloadSize(class int) (int, int) {
	//fetch all attrlen metas
	code := c.codebook.Codes[""]
	var metaAttrLens []int
	for i, meta := range code.Metas {
		if meta.AttrSize {
			metaAttrLens = append(metaAttrLens, i)
		}
	}

	for i, meta := range code.Payloads {
		if meta.AttrSize {
			metaAttrLens = append(metaAttrLens, i)
		}
	}
	return 0, 0
}

func (c *GeneralCoder) getSize(class int) (int, int) {
	return 0, 0
}

func (c *GeneralCoder) decodeMeta(meta *CodeMeta) *MsgAttr {
	return nil
}

func (c *GeneralCoder) decodePayload(meta *CodeMeta) *MsgAttr {
	return nil
}
