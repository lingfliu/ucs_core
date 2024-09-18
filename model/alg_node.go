package model

type AlgNode struct {
	Id    int64
	Class string
	Name  string
	Addr  string
	State int
	//TODO: add 系统负载等
}
