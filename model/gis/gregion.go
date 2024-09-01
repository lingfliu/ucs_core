package gis

type GRegion struct {
	Verts []GPos //顶点
}

/**
 * 计算区域是否封闭，是否存在交叉
 */
func (region *GRegion) EncloseCheck() bool {
	return false
}
