package model

import (
	"github.com/lingfliu/ucs_core/model/meta"
)

// deprecated
type UPoint struct {
	Id         int64
	NodeId     int64
	Name       string
	Descrip    string
	Offset     int //偏置地址，对数据点位，不再区分其地址
	IoMode     int // 1: read(监测), 2: write(control), 3: R & W （双向）
	SampleMode int //read mode: 0: auto, 1: trig, 2: poll
	CtlMode    int //write mode: 0: IO, 1: func
	PropSet    map[string]string
	DataMeta   meta.DataMeta //for read only
	Data       []byte
}
