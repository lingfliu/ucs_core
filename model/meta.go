package model

const (
	DATA_CLASS_RAW   = 0 //in bytes
	DATA_CLASS_INT   = 1
	DATA_CLASS_UINT  = 2
	DATA_CLASS_FLOAT = 3
	DATA_CLASS_FLAG  = 4
	// DATA_CLASS_JSON  = 5 //in json string
)

type PropMeta struct {
	Name      string
	DataClass int
	Dimen     [2]int //Dimen[0]: number of samples, Dimen[1]: number of dimension for each sample
	Sps       int64  //sampling rate in milliseconds, e.g. 50 means 20Hz, 0 means it is a burst event (occur only once)
}
