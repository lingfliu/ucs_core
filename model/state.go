package model

import "github.com/lingfliu/ucs_core/model/gis"

/**
 * 属性数据没有SPS，只有采集时间戳
 * Dimen[0] = 1 is suggested
 */
type State struct {
	//meta
	Meta *PropMeta

	//content
	Ts      int64    //timestamp
	GPos    gis.GPos //全局坐标
	LPos    gis.LPos //局部坐标
	Payload []any    //泛型
}
