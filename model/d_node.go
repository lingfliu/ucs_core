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
	Descrip  string
	DpSet    map[int64]*DPoint
}
