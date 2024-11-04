package msg

import "github.com/lingfliu/ucs_core/model"

/**
 * 寻址模式下的数据消息结构
 */
const (
	CTL_CLASS_IO   = 0
	CTL_CLASS_FUNC = 1
)

/**
 * 针对单个控制点位的控制指令，为方便收发，将控制点位信息直接放入消息中
 */
type CtlMsg struct {
	NodeId   int64  //DD模式下可仅使用ID进行检索，受控端直接注册该ID号
	NodeAddr string // dnode ip 或 URL, 如采用直接访问的方式需要提供
	Offset   int
	Ts       int64
	Idx      int //序号
	Mode     int //消息类型: 0-IO模式，1-函数调用（基于URL）
	Data     []byte
}

/**
 * 控制点位消息转换
 */
func CtlData2Msg(data *model.CtlData) *CtlMsg {
	return &CtlMsg{
		Ts:       cp.Ts,
		Idx:      cp.Idx,
		Mode:     ,
		NodeAddr: cp.NodeAddr,
		NodeId:   cp.NodeId,

	}
}

//TODO: remove it
func CtlMsg2Point(cm *CtlMsg, cp *model.CtlPoint) {
	if cp.NodeId == cm.NodeId && cp.NodeAddr == cm.NodeAddr && cp.Offset == cm.CtlCmd.Offset {
		cp.Ts = cm.Ts
		cp.Idx = cm.Idx
		cp.Data = cm.CtlCmd.Cmd
	}
}
