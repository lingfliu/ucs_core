package model

import "github.com/lingfliu/ucs_core/model/meta"

type UData struct {
	NodeId  int64
	PointId int64
	Addr    string
	Offset  int

	Session string
	Ts      int64
	Idx     int

	DataMeta *meta.DataMeta
	Data     any
}
