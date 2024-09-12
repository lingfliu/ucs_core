package meta

import "encoding/binary"

type DataMeta struct {
	Dimen     int
	ByteLen   int    //1,2,4,8 only
	Alias     string //代号
	Unit      string //单位
	DnClass   int    //节点类型
	DpClass   int    //类型
	DataClass int
	Msb       bool
}

func DataConvert(raw []byte, class int, byteLen int, dimen int, msb bool) []any {
	data := make([]any, dimen)
	switch class {
	case VAL_CLASS_RAW:
		for i, val := range raw {
			data[i] = val
		}
	case VAL_CLASS_INT:
		for i := 0; i < dimen; i++ {
			if byteLen == 4 {
				//int32
				data[i] = int32(binary.BigEndian.Uint32(raw[i*byteLen : (i+1)*byteLen]))
			} else if byteLen == 8 {
				//int64
				data[i] = int64(binary.BigEndian.Uint64(raw[i*byteLen : (i+1)*byteLen]))
			} else {
				panic("unsupported byte length")
			}
		}
	case VAL_CLASS_UINT:
		for i := 0; i < dimen; i++ {
			if byteLen == 4 {
				//int32
				data[i] = uint32(binary.BigEndian.Uint32(raw[i*byteLen : (i+1)*byteLen]))
			} else if byteLen == 8 {
				//int64
				data[i] = uint64(binary.BigEndian.Uint64(raw[i*byteLen : (i+1)*byteLen]))
			} else {
				panic("unsupported byte length")
			}
		}
	case VAL_CLASS_FLOAT:
		for i := 0; i < dimen; i++ {
			if byteLen == 4 {
				//float32
				data[i] = float32(binary.BigEndian.Uint32(raw[i*byteLen : (i+1)*byteLen]))
			} else if byteLen == 8 {
				//float64
				data[i] = float64(binary.BigEndian.Uint64(raw[i*byteLen : (i+1)*byteLen]))
			} else {
				panic("unsupported byte length")
			}
		}
	case VAL_CLASS_FLAG:
		for i := 0; i < dimen; i++ {
			data[i] = raw[i] != 0
		}
	}

	return data
}
