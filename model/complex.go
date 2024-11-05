package model

import "github.com/lingfliu/ucs_core/model/gis"

/**
 * 结构体
 */
type Complex struct {
	Id    int64
	Name  int64
	Pos   *gis.GPos
	Model string //GIS-BIM model

	CamList   []*Cam
	DNodeList []*UNode
}
