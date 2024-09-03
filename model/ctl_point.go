package model

/**
 * 控制点位
 */
type CtlPoint struct {
	Id        int64
	NodeId    int64
	OffsetIdx int
	Name      string
	Meta      *DataMeta
	Data      []any
}
