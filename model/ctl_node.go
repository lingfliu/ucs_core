package model

/**
 * 控制点位
 */
type CtlNode struct {
	Name      string
	Id        int64
	ParentId  int64
	Descrip   string
	Addr      string
	CtlPoints map[int64]*CtlPoint
}
