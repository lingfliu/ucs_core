package model

/**
 * 项目:当前阶段省略了标段定义
 */
type UProject struct {
	Id      int64
	Name    string
	Descrip string
	SiteSet map[int64]*Complex

	//TODO: add 项目组织管理信息
}
