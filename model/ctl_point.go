package model

import "github.com/lingfliu/ucs_core/model/meta"

/**
 * 控制点位, 该结构和DPoint一致
 */
type CtlPoint struct {
	//1. 寻址， 2. 静态属性， 3. 数据格式， 4. 数据
	Id     int64
	Class  string
	Name   string //名称，规格
	Offset int

	NodeId    int64
	NodeName  string
	NodeClass string
	NodeAddr  string

	DataMeta *meta.DataMeta //控制点位只有两种格式：byte和string
	Mode     int
	Data     []byte
}
