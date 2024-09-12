package model

const (
	DNODE_STATE_OK      = 0
	DNODE_STATE_ALARM   = 1
	DNODE_STATE_FAULT   = 2
	DNODE_STATE_OFFLINE = 3
)

/**
 * 监测点位
 */
type DNode struct {
	Id        int64
	ParentId  int64             //归属
	Name      string            //名称
	Class     int               //节点设备类型编码
	Addr      string            //ip or url
	Desc      string            //文字描述，辅助信息
	PropSet   map[string]string //静态属性，string格式
	DPointSet map[int64]*DPoint //数据点位
	State     int               //0-正常，1-告警，2-故障, 3-离线
	//TODO: 是否需要添加多个状态实时数据
}
