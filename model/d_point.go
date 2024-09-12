package model

import (
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/ulog"
)

type DPoint struct {
	Ts       int64
	Id       int64 //dpoint id
	ParentId int64 //dnode id
	Data     []byte
	Meta     *meta.DataMeta
}

func NewDPoint(id, parentId int64, meta *meta.DataMeta, raw_data []byte) *DPoint {
	byteLen := meta.ByteLen
	dimen := meta.Dimen

	if byteLen*dimen != len(raw_data) {
		ulog.Log().E("model", "byte length not match")
		return nil
	} else {
		return &DPoint{
			Id:       id,
			ParentId: parentId,
			Meta:     meta,
			Data:     raw_data,
			Ts:       0,
		}
	}

}
