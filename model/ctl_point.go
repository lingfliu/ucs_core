package model

import "github.com/lingfliu/ucs_core/model/meta"

/**
 * 控制点位, 该结构和DPoint一致
 */
type CtlPoint struct {
	Id       int64
	NodeId   int64
	NodeAddr string
	Offset   int
	Name     string
	Ts       int64
	Idx      int
	Meta     *meta.DataMeta
	Data     []byte
}
