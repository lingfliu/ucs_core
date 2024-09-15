package msg

import "github.com/lingfliu/ucs_core/model/gis"

type PosMsg struct {
	Ts     int64
	Idx    int
	NodeId int64
	GPos   *gis.GPos
	LPos   *gis.LPos
	Velo   *gis.Velo
}
