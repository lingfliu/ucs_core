package model

import (
	"github.com/lingfliu/ucs_core/model/gis"
	"github.com/lingfliu/ucs_core/model/meta"
)

/**
 * Data from a node
 */
type DnData struct {
	Id   int64
	Addr string

	Ts        int64
	Idx       int
	Session   string
	Sps       int64
	SampleLen int
	Pos       *gis.GPos //could be null

	DpDataList []*DpData
}

type DpData struct {
	Offset   int
	DataMeta *meta.DataMeta
	Data     any
}
