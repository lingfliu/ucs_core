package msg

import "github.com/lingfliu/ucs_core/model"

/**
 * 针对同一个控制节点多个控制点位的同步控制指令
 */
type BatchCtlMsg struct {
	Ts          int64
	Idx         int //序号
	Class       int //消息类型: 0-IO模式，1-函数调用（基于URL）
	CtlNodeAddr string
	CtlNodeId   int64
	CtlCmdSet   map[int64]*CtlCmd //该模式下，所有点位均采用字节流，不做编解码
}

func BatchCtlPoint2Msg(cpSet map[int64]*model.CtlPoint) *BatchCtlMsg {
	var ts int64
	var idx int
	var class int
	var nodeAddr string
	var nodeId int64

	for _, cp := range cpSet {
		ts = cp.Ts
		idx = cp.Idx
		class = CTL_CLASS_IO
		nodeAddr = cp.NodeAddr
		nodeId = cp.NodeId
		break
	}

	CtlCmdSet := make(map[int64]*CtlCmd)
	for _, cp := range cpSet {
		CtlCmdSet[cp.Id] = &CtlCmd{
			Offset: cp.Offset,
			Cmd:    cp.Data,
		}

	}
	message := &BatchCtlMsg{
		Ts:          ts,
		Idx:         idx,
		Class:       class,
		CtlNodeAddr: nodeAddr,
		CtlNodeId:   nodeId,
		CtlCmdSet:   CtlCmdSet,
	}

	return message
}
