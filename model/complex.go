package model

/**
 * 结构体
 */
type Complex struct {
	CamSet  map[int64]*Cam
	NodeSet map[int64]*UNode
	MachSet map[int64]*Mach
	//TODO: 算法节点
}
