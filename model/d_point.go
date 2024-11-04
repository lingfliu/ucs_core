package model

import (
	"github.com/lingfliu/ucs_core/model/meta"
)

type DPoint struct {
	//1. 寻址， 2. 静态属性， 3. Node信息, 4. 数据Meta， 5. 数据
	Id     int64
	Class  string //数据点位类别
	Name   string //数据点位名称
	Offset int    //索引 (总线地址)

	NodeId    int64
	NodeName  string //alternative of NodeId for display
	NodeClass string //classes are used to generate stable name
	NodeAddr  string

	Mode int //监测模式(和所述node一致): 0-采样，1-事件触发，2-轮询
	Sps  int64

	Ts      int64
	Idx     int    //序号
	Session string //会话标识

	// data could be values (data array) or url (files)
	DataMeta *meta.DataMeta
	Data     []byte
}
