package model

/**
 * 项目
 */
type UProject struct {
	Id      int64
	Name    string
	Descrip string
	SiteSet map[int64]*Site

	//TODO: add 项目组织管理信息
}
