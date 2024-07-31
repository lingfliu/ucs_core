package coder

type JsonMsg string

// func (msg *Msg) GetAttrVals(name string) []any {
// 	var attr *MsgAttr
// 	if msg.Metas[name] != nil {
// 		attr = msg.Metas[name]
// 	} else {
// 		attr = msg.Payloads[name]
// 	}
// 	return attr.Vals
// }

// func (msg *Msg) ToBytes() string {
// 	bs, err := json.Marshal(msg)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(bs)
// }

const (
	MSG_ATTR_CLASS_INT       = 0
	MSG_ATTR_CLASS_FLOAT     = 1
	MSG_ATTR_CLASS_STR_ASCII = 2
	MSG_ATTR_CLASS_STR_UTF8  = 2
	MSG_ATTR_CLASS_BOOL      = 3
)

type MsgAttrSpec struct {
	Class    int
	Size     int    //>1 if array
	Encoding string //for string
}

type MsgSpec struct {
	Class int //msg class
	//metas
	Meta map[string]*MsgAttrSpec

	//payloads
	Payload map[string]*MsgAttrSpec
}

type Msg interface {
	ToBytes() []byte
	ToJson() string
	GetVal(name string) any
	GetMeta(name string) any
	GetClass() int
}

type UMsg struct {
	Name    string
	Class   int
	Meta    map[string]any
	Payload map[string]any
}

func (msg *UMsg) ToBytes() []byte {
	return make([]byte, 1)
}

func (msg *UMsg) ToJson() string {
	//marshall
	return ""
}

func (msg *UMsg) GetVal(name string) any {
	return msg.Payload[name]
}

func (msg *UMsg) GetMeta(name string) any {
	return msg.Meta[name]
}
