package msg

/**
 * 寻址模式下的数据消息结构
 */
const (
	DP_MSG_CLASS_AUTO = 0
	DP_MSG_CLASS_TRIG = 1
	DP_MSG_CLASS_POLL = 2
)

type DpMsg struct {
	Ts        int64
	Idx       int //序号， 用于辅助判断是否丢包
	Class     int //消息类型: 0-定时采样，1-事件触发，2-轮询
	DNodeId   int64
	DNodeAddr string // dnode 地址
	DpOffset  int    // 数据点位偏移
	Data      []byte //数据如何解析通过查询DNodeId与DpOffset关联的DataMeta进行换算
}
