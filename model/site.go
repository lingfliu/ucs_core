package model

/**
 * 工地
 */
type Site struct {
	CamSet  map[int64]*Cam
	NodeSet map[int64]*UNode
	MachSet map[int64]*Mach
	//TODO: 算法节点
}
