package model

type AlgResult struct {
	Ts        int64 //计算完成时间
	TaskId    int64 //对应的任务ID
	Code      int   //算法码
	Class     int   //
	AlgNodeId int64 //算法节点ID
	Output    []string
}
