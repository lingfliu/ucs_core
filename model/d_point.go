package model

import (
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/ulog"
)

type DPoint struct {
	//1. 寻址， 2. 静态属性， 3. 数据格式， 4. 数据
	Id       int64
	NodeId   int64
	NodeAddr string
	Offset   int
	Name     string //数据点位名称
	Ts       int64
	Idx      int //序号
	// data could be values (data array) or url (files)
	DataMeta *meta.DataMeta
	Data     []byte
}

// TODO: 基于byte、meta、id的装配
func NewDPoint(id int64, nodeId int64, offset int, ts int64, idx int, dataMeta *meta.DataMeta, data []byte) *DPoint {
	byteLen := dataMeta.ByteLen
	dimen := dataMeta.Dimen

	if byteLen*dimen != len(data) {
		ulog.Log().E("model", "byte length not match")
		return nil
	} else {
		return &DPoint{
			Id:       id,
			NodeId:   nodeId,
			Offset:   offset,
			DataMeta: dataMeta,
			Data:     data,
			Ts:       0,
		}
	}
}
