package types

/**
 * 事件、数据流、数据包，均采用该结构进行声明
 */
type Event struct {
	//meta
	Meta *PropMeta

	//content
	Ts      int64 //timestamp
	GPos    GPos  //全局坐标
	LPos    LPos  //局部坐标
	Payload []any //泛型
}
