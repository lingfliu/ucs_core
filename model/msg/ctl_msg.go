package msg

import "github.com/lingfliu/ucs_core/model"

/**
 * 寻址模式下的数据消息结构
 */
const (
	CTL_CLASS_IO   = 0
	CTL_CLASS_FUNC = 1
)

//控制指令数据
type CtlCmd struct {
	Offset int
	Cmd    []byte //原始字节控制指令, FUNC模式下为函数(参数)字符串: ADD(1,2,4)
}

/**
 * 针对单个控制点位的控制指令，为方便收发，将控制点位信息直接放入消息中
 */
type CtlMsg struct {
	NodeId   int64  //DD模式下可仅使用ID进行检索，受控端直接注册该ID号
	NodeAddr string // dnode ip 或 URL, 如采用直接访问的方式需要提供
	Ts       int64
	Idx      int //序号
	Mode     int //消息类型: 0-IO模式，1-函数调用（基于URL）
	CtlCmd   *CtlCmd
}

/**
 * 控制点位消息转换
 */
func CtlPoint2Msg(cp *model.CtlPoint) *CtlMsg {
	return &CtlMsg{
		Ts:       cp.Ts,
		Idx:      cp.Idx,
		Mode:     CTL_CLASS_IO,
		NodeAddr: cp.NodeAddr,
		NodeId:   cp.NodeId,
		CtlCmd: &CtlCmd{
			Offset: cp.Offset,
			Cmd:    cp.Data,
		},
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
