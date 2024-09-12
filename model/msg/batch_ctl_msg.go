package msg

/**
 * 针对同一个控制节点多个控制点位的控制指令
 */
type BatchCtlMsg struct {
	Ts          int64
	Idx         int //序号
	Class       int //消息类型: 0-IO模式，1-函数调用（基于URL）
	CtlNodeAddr string
	CtlNodeId   int64
	CtlDataSet  map[int64]CtlData //该模式下，所有点位均采用字节流，不做编解码
}
