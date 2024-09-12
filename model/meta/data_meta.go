package meta

import "encoding/binary"

const (
	DATA_CLASS_RAW   = 0 //in bytes (string)
	DATA_CLASS_INT   = 1
	DATA_CLASS_UINT  = 2
	DATA_CLASS_FLOAT = 3
	DATA_CLASS_FLAG  = 4
	DATA_CLASS_JSON  = 5 //UTF-8 format json string
)

/**
 * Normally, a data meta declare a single data specification
 */
type DataMeta struct {
	ByteLen   int //1,2,4,8 only
	Dimen     int
	Alias     string //代号
	Unit      string //单位
	DataClass int
	Msb       bool
}

//TODO: remove on release
func DataConvert(raw []byte, class int, byteLen int, dimen int, msb bool) []any {
	data := make([]any, dimen)
	switch class {
	case DATA_CLASS_RAW:
		for i, val := range raw {
			data[i] = val
		}
	case DATA_CLASS_INT:
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
	case DATA_CLASS_UINT:
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
	case DATA_CLASS_FLOAT:
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
	case DATA_CLASS_FLAG:
		for i := 0; i < dimen; i++ {
			data[i] = raw[i] != 0
		}
	}

	return data
}
