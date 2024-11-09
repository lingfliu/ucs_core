package model

const (
	MACH_CLASS_WSP         = 0 //湿喷机
	MACH_CLASS_SIMPLE_VEHI = 1 //简易车辆
)

/**
 * 工程机械 （mobile）
 */
type Mach struct {
	Id          int64
	Name        string
	Class       string // 设备型号
	Addr        string
	PropSet     map[string]string
	DNodeList   []*DNode
	CtlNodeList []*CtlNode
	CamList     []*Cam
}

func (mach *Mach) hasDNode(nodeId int64) bool {
	for _, node := range mach.DNodeList {
		if node.Id == nodeId {
			return true
		}
	}
	return false
}
