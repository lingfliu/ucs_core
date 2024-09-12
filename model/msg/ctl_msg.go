package msg

/**
 * 寻址模式下的数据消息结构
 */
const (
	CTL_CLASS_IO   = 0
	CTL_CLASS_FUNC = 1
)

//控制指令数据
type CtlData struct {
	Offset int
	Data   []byte //原始字节控制指令, FUNC模式下为函数(参数)字符串: ADD(1,2,4)
}

/**
 * 针对单个控制点位的控制指令，为方便收发，将控制点位信息直接放入消息中
 */
type CtlMsg struct {
	Ts          int64
	Idx         int    //序号
	Class       int    //消息类型: 0-IO模式，1-函数调用（基于URL）
	CtlNodeAddr string // dnode ip 或 URL, 如采用直接访问的方式需要提供
	CtlNodeId   int64  //DD模式下可仅使用ID进行检索，受控端直接注册该ID号
	CtlData     *CtlData
}
