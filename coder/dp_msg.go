package coder

type DpMsg struct {
	Ts         int64
	DNodeClass int32
	DNodeId    int64
	ValueList  []int64
}
