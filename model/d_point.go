package model

import (
	"github.com/lingfliu/ucs_core/model/meta"
)

type DPoint struct {
	//1. 寻址， 2. 静态属性， 3. 数据格式， 4. 数据
	Id       int64
	NodeId   int64
	NodeAddr string
	Offset   int
	Name     string //数据点位名称
	Sps      int64  //采样率, 仅在采样模式下有效
	Ts       int64
	Idx      int //序号
	// data could be values (data array) or url (files)
	DataMeta *meta.DataMeta
	Data     []byte
}
