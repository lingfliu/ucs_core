package model

import (
	"encoding/binary"
)

type CtlPoint struct {
	Name      string
	Id        int64
	ParentId  int64
	Descrip   string
	OffsetIdx int
	//attrs for byte conversion
	ValueClass int
	ByteLen    int
	Dimen      int
	SampleLen  int
}

func (cp *CtlPoint) Convert(rawData []byte) any {
	switch cp.ValueClass {
	case DATA_CLASS_INT:
		return int(binary.BigEndian.Uint32(rawData[:4]))
	case DATA_CLASS_FLOAT:
		return float64(binary.BigEndian.Uint64(rawData))
	case DATA_CLASS_BOOL:
		return rawData[0] != 0
	default:
		return nil
	}
}
