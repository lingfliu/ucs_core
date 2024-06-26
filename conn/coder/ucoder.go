package coder

import (
	"bytes"

	"github.com/lingfliu/ucs_core/utils"
)

type UCoder struct {
	buff     *utils.ByteRingBuffer
	Codebook *Codebook

	_buff []byte

	//local variables
	metaLen       int
	msgMetaLen    int
	prePayloadLen int
	decodeState   int
	meta_tmp      []byte
}

func (coder *UCoder) Reset() {
	coder.buff.Flush()
}

func (coder *UCoder) PushDecode(bs []byte, n int) *UMsg {
	foundHeader := false
	coder.buff.Push(bs, n)
	for coder.buff.Capacity > len(coder.Codebook.Header) {
		//peek till the position of msg class meta
		coder.buff.Peek(coder._buff, len(coder.Codebook.Header))
		if bytes.Equal(coder._buff, coder.Codebook.Header) {
			foundHeader = true
		} else {
			coder.buff.Drop(1)
		}
	}

	if !foundHeader {
		return nil
	}

	metaLen := coder.Codebook.CalcMetaByteLen()
	if coder.buff.Capacity < metaLen {
		return nil
	}

	coder.buff.Peek(coder._buff, metaLen)

	//get msg class
	attrSpec := coder.Codebook.MsgClassAttrSpec
	offset := len(coder.Codebook.Header) / 8
	coder.buff.Peek(coder._buff, len(coder.Codebook.Header)+offset+attrSpec.ByteLen)
	msgClass, _ := coder.DecodeAttr(coder._buff, attrSpec).(int)

	msgSpec := coder.Codebook.GetMsgSpec(msgClass)

	payloadLen := coder.Codebook.PreCalcPayloadLen(msgClass)

	if msgSpec.Varlen {
		for _, codeSpec := range msgSpec.MetaList {
			if codeSpec.LenSpec != "" {
				l := coder.DecodeAttr(coder._buff, codeSpec).(int)
				payloadLen += l
			}
		}
		if coder.buff.Capacity < len(coder.Codebook.Header)+metaLen+payloadLen {
			//buff insufficient, skip
			return nil
		}
	}

	//decode meta
	metaSet := make(map[string]any)
	for _, codeSpec := range coder.Codebook.MetaList {
		meta := coder.DecodeAttr(coder._buff, codeSpec)
		metaSet[codeSpec.Name] = meta
	}

	//msg specified meta
	for _, codeSpec := range msgSpec.MetaList {
		meta := coder.DecodeAttr(coder._buff, codeSpec)
		metaSet[codeSpec.Name] = meta
	}

	payloadSet := make(map[string]any)
	for _, codeSpec := range msgSpec.PayloadList {
		payload := coder.DecodeAttr(coder._buff, codeSpec)
		payloadSet[codeSpec.Name] = payload
	}

	return &UMsg{
		Name:    msgSpec.Name,
		Class:   msgSpec.Class,
		Meta:    metaSet,
		Payload: payloadSet,
	}
}

func (coder *UCoder) DecodeAttr(bs []byte, attrSpec *CodeAttrSpec) any {
	switch attrSpec.Class {
	case CODE_CLASS_INT:
		values := make([]int, attrSpec.Size)
		offset := 0
		for i := 0; i < attrSpec.Size; i++ {
			values[i] = utils.Byte2Int(bs, offset, attrSpec.ByteLen, attrSpec.Unsigned, attrSpec.Msb)
		}
		return values
	case CODE_CLASS_FLOAT:
		offset := 0
		if attrSpec.ByteLen == 4 {
			values := make([]float32, attrSpec.Size)

			for i := 0; i < attrSpec.Size; i++ {
				values[i] = utils.Byte2Float(bs, offset, attrSpec.Msb)
				offset += 4
			}
			return values
		} else if attrSpec.ByteLen == 8 {
			values := make([]float64, attrSpec.Size)
			for i := 0; i < attrSpec.Size; i++ {
				values[i] = utils.Byte2Double(bs, offset, attrSpec.Msb)
				offset += 8
			}
			return values
		} else {
			return nil
		}

	case CODE_CLASS_STRING:
		return utils.Byte2String(bs, attrSpec.Offset, attrSpec.Size)

		// case CODE_CLASS_FLAG:
		// 	values := make([]bool, attrSpec.Size)
		// 	offset := 0
		// 	for i := 0; i < attrSpec.Size; i++ {
		// 		values[i] = utils.Byte2Bool(bs, offset)
		// 		offset += attrSpec.ByteLen
		// 	}
	default:
		return nil
	}
}

func (coder *UCoder) Encode(msg *UMsg, bs []byte) int {
	msgCodeSpec := coder.Codebook.GetMsgSpec(msg.Class)
	// byteLen := coder.Codebook.CalcMsgByteLen(msg.Class)

	idx := 0
	//copy header
	lheader := len(coder.Codebook.Header)
	copy(bs[idx:lheader], coder.Codebook.Header)

	idx += lheader
	//put meta
	for _, codeSpec := range msgCodeSpec.MetaList {
		coder.EncodeAttr(msg.Meta[codeSpec.Name], codeSpec, bs)
	}

	//put payload
	for _, codeSpec := range msgCodeSpec.PayloadList {
		coder.EncodeAttr(msg.Payload[codeSpec.Name], codeSpec, bs)
	}

	return 0
}

func (coder *UCoder) EncodeAttr(value any, spec *CodeAttrSpec, bs []byte) int {
	return 0
}
