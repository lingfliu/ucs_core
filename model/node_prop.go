package model

/**
 * 静态属性数据没有时间戳
 */
type NodeProp struct {
	Name string
	Data []any
	Meta *DataMeta
}
