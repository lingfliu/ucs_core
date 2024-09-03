package model

import "github.com/lingfliu/ucs_core/model/gis"

/**
 * 预警事件，均采用该结构进行声明
 */
type Event struct {
	Ts   int64
	GPos gis.GPos //全局坐标
	LPos gis.LPos //局部坐标

	Code       int       //事件码
	Snapshot   []*DPoint //数据快照
	AlertLevel int       //告警级别
}
