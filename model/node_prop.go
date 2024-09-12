package model

import (
	"encoding/binary"

	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/utils"
)

/**
 * 静态属性数据没有时间戳
 */
type NodeProp struct {
	Name string
	Data []byte
	Meta *meta.PropMeta
}

func (p *NodeProp) ToString() string {
	bl := p.Meta.ByteLen
	di := p.Meta.Dimen
	if bl == 1 {
		//convert to string
		return string(p.Data)
	} else if bl == 2 {
		//uint16
		dlist := make([]uint16, di)
		for i := 0; i < di; i++ {
			dlist[i] = binary.BigEndian.Uint16(p.Data[i*bl : (i+1)*bl])
		}
	} else if bl == 4 {
		//uint32
		dlist := make([]uint32, di)
		for i := 0; i < di; i++ {
			dlist[i] = binary.BigEndian.Uint32(p.Data[i*bl : (i+1)*bl])
		}
		str := ""
		for _, v := range dlist {
			str += string(v)
			str += " "
		}
		return str[:len(str)-1]

	} else if bl == 8 {
		//uint64
		dlist := make([]any, di)
		for i := 0; i < di; i++ {
			dlist[i] = binary.BigEndian.Uint64(p.Data[i*bl : (i+1)*bl])
		}

		utils.Array2String(dlist)
		str := ""
		for _, v := range dlist {
			str += string(v)
			str += " "
		}
		return str[:len(str)-1]

	} else {
		//unsupported
		return ""
	}
}
