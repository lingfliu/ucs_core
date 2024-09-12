package model

import (
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/ulog"
)

type DPoint struct {
	Ts     int64
	Idx    int   //序号
	Id     int64 //dpoint id
	NodeId int64 //dnode id
	Name   string
	Offset int //index offset
	/*
	 * data could be values (data array) or url (files)
	 */
	Data []byte
	Meta *meta.DataMeta
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
			Id:     id,
			NodeId: nodeId,
			Offset: offset,
			Meta:   dataMeta,
			Data:   data,
			Ts:     0,
		}
	}
}
