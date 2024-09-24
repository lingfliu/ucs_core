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
	Ts       int64
	Idx      int    //序号
	Session  string //会话标识
	// data could be values (data array) or url (files)
	DataMeta *meta.DataMeta
	Data     []byte
}
