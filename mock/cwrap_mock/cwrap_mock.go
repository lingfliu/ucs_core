package main

import (
	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/ulog"
)

func main() {
	ulog.Config(ulog.LOG_LEVEL_INFO, "", false)

	pt := &model.DPoint{
		DataMeta: &meta.DataMeta{
			ByteLen:   4,
			Dimen:     3,
			SampleLen: 1,
		},
	}
	rtdb.CreatePtTable(pt)

	rtdb.InserDData(pt)
}
