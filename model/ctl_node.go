package model

/**
 * 控制节点
 * 控制节点包括两种模式：
 * 0. IO模式，如对PLC寄存器的写入
 * 1. Function模式，如对设备服务的调用，在该模式下，Addr包含了服务名称和参数，如"SetLedColor(0, 255, 0)", CtlPoints为空
 * 默认情况下CtlNode的数据均为一次控制
 */
type CtlNode struct {
	Id           int64
	ParentId     int64             //mach / complex id
	Addr         string            //url / ip
	Name         string            //名称 / 编号
	Class        string            //类型
	Mode         int               // 控制模式 0-IO模式，1-Function模式
	Descrip      string            //(optional)
	PropSet      map[string]string //静态属性集(optional)
	CtlPointList []*CtlPoint
	State        int //状态 0-正常， 2-故障， 3-离线
}
