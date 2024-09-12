package model

/**
 * 算法任务
 */
type AlgTask struct {
	Ts       int64 //提交时间
	Id       int64
	Class    string //算法类型
	Input    []string
	Priority int //计算优先级
}
