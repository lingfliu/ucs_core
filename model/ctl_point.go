package model

type CtlPoint struct {
	Id        int64
	ParentId  int64
	OffsetIdx int
	Meta      *DataMeta
}
