package model

const (
	DNODE_STATE_OK      = 0
	DNODE_STATE_ALARM   = 1
	DNODE_STATE_FAULT   = 2
	DNODE_STATE_OFFLINE = 3

	DNODE_MODE_AUTO = 0
	DNODE_MODE_TRIG = 1
	DNODE_MODE_POLL = 2
)

/**
 * 监测点位
 */
type DNode struct {
	Id        int64
	ParentId  int64             //归属
	Addr      string            //url
	Name      string            //名称
	Class     string            //节点设备类型编码
	Mode      int               //监测模式: 0-采样，1-事件触发，2-轮询
	Descrip   string            //文字描述，辅助信息
	PropSet   map[string]string //静态属性，string格式
	DPointSet map[int64]*DPoint //数据点位
	State     int               //0-正常，1-告警，2-故障, 3-离线
	//其他的state定义为若干个DPoint
}
