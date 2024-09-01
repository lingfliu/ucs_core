package model

import "github.com/lingfliu/ucs_core/model/gis"

/**
 * 事件、数据流、数据包，均采用该结构进行声明
 */
type Event struct {
	//meta
	Meta *PropMeta

	//content
	Ts      int64    //timestamp
	GPos    gis.GPos //全局坐标
	LPos    gis.LPos //局部坐标
	Payload []any    //泛型
}
