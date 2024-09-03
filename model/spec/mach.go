package spec

import "github.com/lingfliu/ucs_core/model"

/**
 * 大型工程设备
 */
type Mach struct {
	model.Thing

	DNodeSet   map[string]*model.DNode
	CtlNodeSet map[string]*model.CtlNode
}
