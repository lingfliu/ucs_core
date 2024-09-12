package model

import "github.com/lingfliu/ucs_core/model/meta"

/**
 * 控制点位
 */
type CtlPoint struct {
	Id     int64
	NodeId int64
	Offset int
	Name   string
	Meta   *meta.DataMeta
	Data   []byte
}
