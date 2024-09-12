package model

/**
 * 算法结果
 */
type AlgResult struct {
	Ts          int64             //计算完成时间
	TaskId      int64             //对应的任务ID
	Class       string            //算法类型
	Code        int               //状态码
	AlgNodeId   int64             //算法节点ID
	AlgNodeProp map[string]string //算法节点属性
	Output      string            //json格式的输出
}
