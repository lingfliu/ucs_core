package msg

/**
 * 寻址模式下的数据消息结构
 */
const (
	CTL_CLASS_IO   = 0 //CMD as raw bytes
	CTL_CLASS_FUNC = 1 //CMD as FUNC(32bytes) + ARGS_NUM(1) + ARGS(8*ARGS_NUM)
)

type CtlCmd struct {
	Offset int
	Data   []byte
}

/**
 * 基于URL的控制消息
 */
type CtlMsg struct {
	NodeId   int64  //DD模式下可仅使用ID进行检索，受控端直接注册该ID号
	NodeAddr string // dnode ip 或 URL, 如采用直接访问的方式需要提供
	Ts       int64
	Idx      int    //序号
	Session  string //optional

	Mode   int //控制类型: 0-IO模式，1-Function模式
	CtlCmd []*CtlCmd
}

// /**
//  * 控制点位消息转换
//  */
// func CtlData2Msg(data *model.CtlData) *CtlMsg {
// 	return &CtlMsg{
// 		Ts:       cp.Ts,
// 		Idx:      cp.Idx,
// 		Mode:     ,
// 		NodeAddr: cp.NodeAddr,
// 		NodeId:   cp.NodeId,

// 	}
// }

// //TODO: remove it
// func CtlMsg2Point(cm *CtlMsg, cp *model.CtlPoint) {
// 	if cp.NodeId == cm.NodeId && cp.NodeAddr == cm.NodeAddr && cp.Offset == cm.CtlCmd.Offset {
// 		cp.Ts = cm.Ts
// 		cp.Idx = cm.Idx
// 		cp.Data = cm.CtlCmd.Cmd
// 	}
// }
