package model

/**
 * 监测点位
 */
type DNode struct {
	Id       int64
	ParentId int64
	Name     string
	Class    int
	Addr     string //ip or url
	Desc     string //文字描述，辅助信息
	DpSet    map[int64]*DPoint
}
