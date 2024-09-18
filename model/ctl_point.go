package model

import "github.com/lingfliu/ucs_core/model/meta"

/**
 * 控制点位, 该结构和DPoint一致
 */
type CtlPoint struct {
	//1. 寻址， 2. 静态属性， 3. 数据格式， 4. 数据
	Id       int64
	NodeId   int64
	NodeAddr string
	Offset   int
	Name     string //名称，规格
	Ts       int64
	Idx      int
	DataMeta *meta.DataMeta
	Data     []byte
}
