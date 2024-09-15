package model

import "github.com/lingfliu/ucs_core/model/gis"

/**
 * 预警事件
 */
type AlertEvent struct {
	Ts   int64
	GPos gis.GPos //全局坐标
	LPos gis.LPos //局部坐标

	Desc     string    //描述
	Code     int       //事件码
	Level    int       //告警级别
	Snapshot []*DPoint //数据快照
}
