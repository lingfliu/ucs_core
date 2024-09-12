package model

/**
 * 大型工程设备
 */
const (
	MACH_CLASS_WSP = 0 //湿喷机
)

type Mach struct {
	Id         int64
	Name       string
	Addr       string
	Class      int //设备型号
	PropSet    map[string]string
	DNodeSet   map[int64]*DNode   //id as key
	CtlNodeSet map[int64]*CtlNode //id as key
	CamSet     map[int64]*Cam     //id as key
}
