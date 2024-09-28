package msg

/**
 * 寻址模式下的数据消息结构
 */
const (
	D_MSG_CLASS_AUTO = 0
	D_MSG_CLASS_TRIG = 1
	D_MSG_CLASS_POLL = 2
)

type DData struct {
	Offset int
	Data   []byte
}

/**
 * 监测点位消息
 */
type DMsg struct {
	DNodeId   int64
	DNodeAddr string
	Offset    int
	Ts        int64
	Idx       int            //序号， 用于辅助判断是否丢包
	Session   string         //会话标识
	Mode      int            //模式，对应DNode的Mode: 0-定时采样，1-事件触发，2-轮询
	Sps       int64          //采样频率, 仅在Mode=0时有效
	DataSet   map[int]*DData //offset as the key
}
