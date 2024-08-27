package coder

import (
	"encoding/json"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

type Codebook struct {
	Header           []byte
	Tail             []byte
	HasTail          bool
	HasPayload       bool
	MetaList         []*CodeAttrSpec
	MsgSet           map[int]*CodeMsgSpec
	Checksum         string
	MsgClassAttrSpec *CodeAttrSpec
}

func NewCodebookFromJson(jsonStr string) *Codebook {
	var codebook *Codebook
	err := json.Unmarshal([]byte(jsonStr), &codebook)
	if err != nil {
		return nil
	}
	if codebook == nil {
		ulog.Log().I("codebook", "json unmarshal failed")
		return nil
	}
	if codebook.Validate() != "passed" {
		ulog.Log().I("codebook", "validation failed")
		return nil
	} else {
		return codebook
	}
}

func (cb *Codebook) GetMetaSpec(name string) *CodeAttrSpec {
	for _, spec := range cb.MetaList {
		if spec.Name == name {
			return spec
		}
	}
	return nil
}

func (cb *Codebook) GetMsgSpec(class int) *CodeMsgSpec {
	if spec, ok := cb.MsgSet[class]; ok {
		return spec
	} else {
		return nil
	}
}

/**
 * calculate mss byte length given meta bytes
 * all varlen attrs should declare their lengths in meta
 */
func (cb *Codebook) CalcMsgByteLen(msgClass int, meta []byte) int {
	codeMsgSpec := cb.GetMsgSpec(msgClass)
	cnt := len(cb.Header)
	cnt += cb.CalcMetaByteLen()
	for _, codeSpec := range codeMsgSpec.MetaList {
		if codeSpec.LenSpec != "" {
			cnt += utils.Byte2Int(meta, codeSpec.Offset, codeSpec.ByteLen, false, codeSpec.Msb)
		}
	}

	for _, codeSpec := range codeMsgSpec.PayloadList {
		if codeSpec.Size > 0 {
			cnt += codeSpec.ByteLen * codeSpec.Size
		}
	}
	return cnt
}

func (cb *Codebook) CalcMetaByteLen() int {
	meta := cb.MetaList[len(cb.MetaList)-1]
	return meta.Offset + meta.ByteLen*meta.Size
}

func (cb *Codebook) PreCalcPayloadLen(msgClass int) int {
	codeMsgSpec := cb.GetMsgSpec(msgClass)
	cnt := 0
	for _, codeSpec := range codeMsgSpec.PayloadList {
		if codeSpec.Size > 0 {
			cnt += codeSpec.ByteLen * codeSpec.Size
		}
	}
	return cnt
}

/**
 * Validate the codebook
 * return error message if failed
 */
func (cb *Codebook) Validate() string {
	//TODO: implement func
	return "passed"
}
