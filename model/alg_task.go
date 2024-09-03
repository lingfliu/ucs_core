package model

type AlgTask struct {
	Ts       int64
	Id       int64
	Class    int
	Input    []string
	Priority int //计算优先级
}
