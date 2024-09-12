package msg

/**
 * 寻址模式下的数据消息结构
 */
const (
	CTL_MSG_CLASS_IO   = 0
	CTL_MSG_CLASS_FUNC = 1
)

type CtlMsg struct {
	Ts          int64
	Idx         int    //序号
	Class       int    //消息类型: 0-IO模式，1-函数调用（基于URL）
	CtlNodeAddr string // dnode ip 或 URL
	CtlNodeId   int64
	CpOffset    int    // 控制点位偏移
	Data        []byte //原始字节控制指令, FUNC模式下为函数(参数): ADD(1,2,4)
}
