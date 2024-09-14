package model

/**
 * 控制点位
 * 控制点位包括两种模式：
 * 0. IO模式，如对PLC寄存器的写入
 * 1. Function模式，如对设备服务的调用，在该模式下，Addr包含了服务名称和参数，如"SetLedColor(0, 255, 0)", CtlPoints为空
 */
type CtlNode struct {
	Id        int64
	ParentId  int64
	Addr      string
	Name      string
	Class     int // 0-IO模式，1-Function模式
	Desc      string
	PropSet   map[string]string //静态属性集
	CtlPoints map[int64]*CtlPoint
	State     int //状态 0-正常， 2-故障， 3-离线
}
