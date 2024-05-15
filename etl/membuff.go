package etl

import (
	"time"

	"github.com/lingfliu/ucs_core/types"
)

type MemBuff struct {
	Buff     map[string][][]any //maximal support for 2D data [time][data + pos]
	BuffSize map[string]int     //time dimension for different data
	BuffTs   map[string]int64

	TsBaseline int64
	TsIdx      int64
	TsStep     int64
}

func NewMemBuff(TsStep int64) *MemBuff {
	mb := &MemBuff{
		TsStep: TsStep,
	}
	mb.Init()
	return mb
}

func (mb *MemBuff) Init() {
	mb.Buff = make(map[string][][]any)
	mb.BuffSize = make(map[string]int)

	//create Ts baseline
	mb.TsBaseline = time.Now().UnixNano()
	tsBaseline, tsIdx := calcTsBaseline(mb.TsStep)

	mb.TsBaseline = tsBaseline
	mb.TsIdx = tsIdx
}

func (mb *MemBuff) Push(name string, data *types.Data) {
	// stData = mb.StSlice()
}

// func (mb *MemBuff) StSlice(step int64, data types.Data) *types.Data {

// 	tsIdx := (data.Ts - mb.TsBaseline) % mb.TsStep

// 	dBaseline := &types.Data{
// 		Meta:    data.Meta,
// 		Ts:      mb.TsBaseline + tsIdx*mb.TsStep,
// 		Payload: data.Payload,
// 	}

// 	return dBaseline
// }

/**
 * Merge data from different data source with consistent time idx
 * mode:
 */
func (mb *MemBuff) Merge(mode int, dataName ...string) *types.Data {
	// dataMerge := types.Data{
	// 	Meta: &types.PropMeta{
	// 		Name: "Merge"},
	// 	Ts: mb.BuffTs[dataName[0]],
	// }

	// for _, name := range dataName {
	// 	dataBuff := mb.Buff[name]
	// 	dataTs := mb.BuffTs[name]
	// }
	// stData := make(map[string]*types.Data)
	// dimenData := make(map[string]int)
	// for _, name := range dataName {
	// dimenData[name] = mb.Buff[name][0][0].(int)
	// }

	data := types.Data{
		Meta: &types.PropMeta{
			Name: "Merge",
		},
		Ts:      mb.BuffTs[dataName[0]],
		Payload: make([]any, 0),
	}

	// for _, name := range dataName {
	// stData[name] = mb.StSlice(name)
	// data.Payload[0] = append(data.Payload[0], stData[name].Payload[0])
	// }

	return &data
}

func (mb *MemBuff) StSlice(dataName string, timeStep int64, align bool) *types.Data {
	ts := mb.BuffTs[dataName]
	if align {
		ts = (ts%timeStep + 1) * ts
	}

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

	return &types.Data{}
}

func calcTsBaseline(tsStep int64) (int64, int64) {
	ts := time.Now().UnixNano()
	idx := ts % tsStep

	tsBaseline := idx * tsStep

	return tsBaseline, idx
}
