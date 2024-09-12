package spec

import "github.com/lingfliu/ucs_core/model"

/**
 * 湿喷机
 * 固定规格:
 */
func NewWspMach(id int64, name string, addr string, desc string, propSet map[string]string, DNodeSet map[int64]*model.DNode, CtlNodeSet map[int64]*model.CtlNode, CamSet map[int64]*model.Cam) *model.Mach {
	mach := &model.Mach{
		Id:         id,
		Name:       name,
		Addr:       addr,
		Class:      model.MACH_CLASS_WSP,
		PropSet:    propSet,
		DNodeSet:   DNodeSet,
		CtlNodeSet: CtlNodeSet,
		CamSet:     CamSet,
	}
	return mach
}
