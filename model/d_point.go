package model

import (
	"github.com/lingfliu/ucs_core/model/meta"
)

type DPoint struct {
	//1. 寻址， 2. 静态属性， 3. Node信息, 4. 数据
	Id     int64
	Class  string //数据点位类别
	Alias  string //数据点位别称
	Offset int    //索引 (总线地址)

	NodeId    int64
	NodeName  string //alternative of NodeId for display
	NodeClass string //classes are used to generate stable name
	NodeAddr  string

	DataMeta *meta.DataMeta
	// data could be values (data array) or url (files)
	Data []byte //实时监测值
}
