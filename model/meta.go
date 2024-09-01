package model

const DATA_CLASS_INT = 1
const DATA_CLASS_FLOAT = 2
const DATA_CLASS_BOOL = 3

type PropMeta struct {
	Name      string
	DataClass int
	Dimen     [2]int //Dimen[0]: number of samples, Dimen[1]: number of dimension for each sample
	Sps       int64  //sampling rate in milliseconds, e.g. 50 means 20Hz, 0 means it is a burst event (occur only once)
}
