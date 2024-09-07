package buff

import (
	"math"
	"time"

	"github.com/lingfliu/ucs_core/model/gis"
)

type MemData struct {
	Data []uint64 //默认使用uint64进行存储
	Ts   int64
	Pos  gis.LPos
}

func (md *MemData) ReturnAsFloat64() []float64 {
	data := make([]float64, len(md.Data))
	for i, d := range md.Data {
		//from byte to float64
		data[i] = math.Float64frombits(d)
	}
	return data
}

func (md *MemData) ReturnAsUint64() []uint64 {
	return md.Data
}

func (md *MemData) ReturnAsUint32() []uint32 {
	data := make([]uint32, len(md.Data))
	for i, d := range md.Data {
		data[i] = uint32(d)
	}
	return data
}

type MemBuff struct {
	Buff   map[string][]MemData //maximal support for 2D data [ts][pos index]
	BuffTs map[string][]int64   //time index
}

func NewMemBuff(maxTsLen int, maxSpLen int) *MemBuff {
	mb := &MemBuff{}
	return mb
}

func (mb *MemBuff) Push(name string, data []MemData) {
	mb.Buff[name] = append(mb.Buff[name], data...)
}

/**
 * Spatial-Temporal buff slice
 * @param name: name of the buffer
 * @param tic: start time
 */
func (mb *MemBuff) StSlice(name string, tic int64, toc int64, tStep int64, sStep float32, sparse bool) *MemData {

	//create time index
	// tsScale := make([]int64, 10)
	// for i, t := range tsScale {
	// 	tsScale[i] = ts + int64(i)*timeStep
	// }

	// tsScale = []int64{ts}
	// t0 = ts
	// for range(0, )
	// 	tsScale = append(tsScale, t0)
	// 	t0 += timeStep
	// }

	return &MemData{}
}

func CalcTsBaseline(tsStep int64) (int64, int64) {
	ts := time.Now().UnixNano()
	idx := ts % tsStep

	tsBaseline := idx * tsStep

	return tsBaseline, idx
}

func (mb *MemBuff) Reg(name string, buffdimen int, bufflen int) int {
	//if name is in the keys of Buff
	if _, ok := mb.Buff[name]; ok {
		// mb.Buff[name] = make([]float64, bufflen)
		return 0
	} else {
		//duplicate reg, return err
		return -1
	}
}

func (mb *MemBuff) Unreg(name string) {
	delete(mb.Buff, name)
}
