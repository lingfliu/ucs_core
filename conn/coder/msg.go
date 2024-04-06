package coder

import "encoding/json"

type Msg struct {
	Class     int
	Metas     map[string]*MsgAttr
	Payloads  map[string]*MsgAttr
	Timestamp int64
}

type MsgAttr struct {
	Name string

	//container for different types of value
	ValClass int
	Vals     []any
}

func (msg *Msg) GetAttrVals(name string) []any {
	var attr *MsgAttr
	if msg.Metas[name] != nil {
		attr = msg.Metas[name]
	} else {
		attr = msg.Payloads[name]
	}
	return attr.Vals
}

func (msg *Msg) ToString() string {
	bs, err := json.Marshal(msg)
	if err != nil {
		return ""
	}

	return string(bs)
}
