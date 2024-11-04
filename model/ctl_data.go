package model

import (
	"github.com/lingfliu/ucs_core/model/meta"
)

/**
 * Ctl cmd to a ctlPoint
 */
type CtlData struct {
	NodeId   int64
	NodeAddr string
	Offset   int

	Ts      int64  //optional
	Idx     int    //optional
	Session string //optional

	Mode     int //inherit from ctlpoint
	DataMeta *meta.DataMeta
	Data     any
}
