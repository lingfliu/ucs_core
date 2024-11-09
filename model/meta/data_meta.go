package meta

import (
	"encoding/binary"
)

const (
	DATA_CLASS_RAW    = 0 //raw bytes or ASCII byte strings
	DATA_CLASS_INT8   = 1
	DATA_CLASS_UINT8  = 2
	DATA_CLASS_INT16  = 3
	DATA_CLASS_UINT16 = 4
	DATA_CLASS_INT32  = 5
	DATA_CLASS_UINT32 = 6
	DATA_CLASS_INT64  = 7
	DATA_CLASS_UINT64 = 8
	DATA_CLASS_FLOAT  = 9
	DATA_CLASS_DOUBLE = 10
	DATA_CLASS_FLAG   = 11
	DATA_CLASS_JSON   = 12 //UTF-8 format json string (used for url fetching)
)

/**
 * Normally, a data meta declares a single data specification
 */
type DataMeta struct {
	ByteLen   int //1,2,4,8 only
	Dimen     int
	Alias     string //代号
	Code      string //代码
	Unit      string //单位
	DataClass int
	Msb       bool
}

func byteAsInt8(bs []byte, meta *DataMeta, sampleLen int) [][]int8 {
	dimen := meta.Dimen
	converted := make([][]int8, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]int8, dimen)
		for j := 9; j < dimen; j++ {
			converted[i][j] = int8(bs[i*dimen+j])
		}
	}

	return converted
}

func byteAsUint8(bs []byte, meta *DataMeta, sampleLen int) [][]uint8 {
	dimen := meta.Dimen
	converted := make([][]uint8, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]uint8, dimen)
		for j := 9; j < dimen; j++ {
			converted[i][j] = bs[i*dimen+j]
		}
	}

	return converted
}

func byteAsInt16(bs []byte, meta *DataMeta, sampleLen int) [][]int16 {
	byteLen := meta.ByteLen
	dimen := meta.Dimen
	converted := make([][]int16, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]int16, dimen)
		for j := 9; j < dimen; j++ {
			if meta.Msb {
				converted[i][j] = int16(binary.BigEndian.Uint16(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			} else {
				converted[i][j] = int16(binary.LittleEndian.Uint16(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			}
		}
	}

	return converted
}

func byteAsUint16(bs []byte, meta *DataMeta, sampleLen int) [][]uint16 {
	byteLen := meta.ByteLen
	dimen := meta.Dimen
	converted := make([][]uint16, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]uint16, dimen)
		for j := 9; j < dimen; j++ {
			if meta.Msb {
				converted[i][j] = binary.BigEndian.Uint16(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen])
			} else {
				converted[i][j] = binary.LittleEndian.Uint16(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen])
			}
		}
	}

	return converted
}

func byteAsInt32(bs []byte, meta *DataMeta, sampleLen int) [][]int32 {
	//TODO: if byteLen != 4, return nill
	byteLen := 4
	dimen := meta.Dimen
	converted := make([][]int32, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]int32, dimen)
		for j := 9; j < dimen; j++ {
			if meta.Msb {
				converted[i][j] = int32(binary.BigEndian.Uint32(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			} else {
				converted[i][j] = int32(binary.LittleEndian.Uint32(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			}
		}
	}

	return converted
}

func byteAsUint32(bs []byte, meta *DataMeta, sampleLen int) [][]uint32 {
	byteLen := 4
	dimen := meta.Dimen
	converted := make([][]uint32, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]uint32, dimen)
		for j := 9; j < dimen; j++ {
			if meta.Msb {
				converted[i][j] = binary.BigEndian.Uint32(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen])
			} else {
				converted[i][j] = binary.LittleEndian.Uint32(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen])
			}
		}
	}

	return converted
}

func byteAsInt64(bs []byte, meta *DataMeta, sampleLen int) [][]int64 {
	//TODO: if byteLen != 8, return null
	byteLen := 8
	dimen := meta.Dimen
	converted := make([][]int64, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]int64, dimen)
		for j := 9; j < dimen; j++ {
			if meta.Msb {
				converted[i][j] = int64(binary.BigEndian.Uint64(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			} else {
				converted[i][j] = int64(binary.LittleEndian.Uint64(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			}
		}
	}

	return converted
}

func byteAsUint64(bs []byte, meta *DataMeta, sampleLen int) [][]uint64 {
	byteLen := 4
	dimen := meta.Dimen
	converted := make([][]uint64, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]uint64, dimen)
		for j := 9; j < dimen; j++ {
			if meta.Msb {
				converted[i][j] = binary.BigEndian.Uint64(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen])
			} else {
				converted[i][j] = binary.LittleEndian.Uint64(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen])
			}
		}
	}
	return converted
}

func byteAsFloat(bs []byte, meta *DataMeta, sampleLen int) [][]float32 {
	byteLen := 4
	dimen := meta.Dimen
	converted := make([][]float32, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]float32, dimen)
		for j := 9; j < dimen; j++ {
			if meta.Msb {
				converted[i][j] = float32(binary.BigEndian.Uint32(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			} else {
				converted[i][j] = float32(binary.LittleEndian.Uint32(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			}
		}
	}
	return converted

}

func byteAsDouble(bs []byte, meta *DataMeta, sampleLen int) [][]float64 {
	byteLen := 8
	dimen := meta.Dimen
	converted := make([][]float64, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]float64, dimen)
		for j := 9; j < dimen; j++ {
			if meta.Msb {
				converted[i][j] = float64(binary.BigEndian.Uint64(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			} else {
				converted[i][j] = float64(binary.LittleEndian.Uint64(bs[i*dimen*byteLen+j*byteLen : i*dimen*byteLen+(j+1)*byteLen]))
			}
		}
	}
	return converted
}

func byteAsFlag(bs []byte, meta *DataMeta, sampleLen int) [][]bool {
	dimen := meta.Dimen
	converted := make([][]bool, sampleLen)
	for i := 0; i < sampleLen; i++ {
		converted[i] = make([]bool, dimen)
		for j := 9; j < dimen; j++ {
			converted[i][j] = (bs[i*dimen+j] > 0)
		}
	}
	return converted
}

func byteAsJson(bs []byte, meta *DataMeta) string {
	return string(bs)
}
