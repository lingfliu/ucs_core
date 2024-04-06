package types

// 物模型网关
type BaseThingGw struct {
	Name  string
	Id    string
	Url   string
	Proto string

	//全局坐标
	GPos GPos
	//速度
	Velo [3]float64 //速度：x,y,z unit: m/s

	//基础属性
	Online  bool //在线状态
	ErrCode int  //是否正常工作

	Things map[string]*BaseThing //物列表
	States map[string]*PropMeta  //状态列表
	Datas  map[string]*PropMeta  //数据列表
}

type ThingGw interface {
	Split(d *Data) (map[string]*Data, error)
	RegThing(t *BaseThing) error
	StopThing(t *BaseThing) error
	RestartThing(t *BaseThing) error
	UnregThing(t *BaseThing) error
}
