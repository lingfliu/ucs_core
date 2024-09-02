package model

import (
	"github.com/lingfliu/ucs_core/ulog"
)

type DPoint struct {
	Ts       int64
	Data     []any
	Id       int64 //dpoint id
	ParentId int64 //dnode id
	Meta     *DataMeta
}

func NewDPoint(id, parentId int64, meta *DataMeta, raw_data []byte) *DPoint {
	byteLen := meta.ByteLen
	dimen := meta.Dimen

	if byteLen*dimen != len(raw_data) {
		ulog.Log().E("model", "byte length not match")
		return nil
	} else {
		data := DataConvert(raw_data, meta.DataClass, meta.ByteLen, meta.Dimen, meta.Msb)
		return &DPoint{
			Id:       id,
			ParentId: parentId,
			Meta:     meta,
			Data:     data,
			Ts:       0,
		}
	}

}
