package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/lingfliu/ucs_core/dao"
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/model/msg"
	"github.com/lingfliu/ucs_core/ulog"
	_ "github.com/taosdata/driver-go/v3/taosSql"
)

const (
	TAOS_HOST     = "62.234.16.239:6030"
	TAOS_DATABASE = "ucs"
	TAOS_USERNAME = "root"
	TAOS_PASSWORD = "taosdata"
)

func main() {
	ulog.Config(ulog.LOG_LEVEL_DEBUG, "", false)

	//config taos
	//TODO: fix the password
	dpDao := dao.NewDpDao(TAOS_HOST, TAOS_DATABASE, TAOS_USERNAME, TAOS_PASSWORD)
	go _task_dao_init(dpDao)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	for {
		select {
		case <-s:
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func _task_dao_query(dao *dao.DpDao) {
	tic := "2024-01-01 00:00:00.000"
	toc := "2024-11-10 00:00:00.000"

	ptList := dao.Query(tic, toc, 1, 1, &meta.DataMeta{
		Dimen:   4,
		ByteLen: 4,
	})
	for _, pt := range ptList {
		//serialize pt
		ulog.Log().I("main", string(pt.Data))
	}
}

func _task_insert(dao *dao.DpDao) {
	dmsg := &msg.DMsg{
		DNodeId: 1,
		Offset:  0,
		Ts:      time.Now().UnixNano() / 1000000,
	}

	dmsg.DataList = make(map[int]*msg.DMsgData)
	dmsg.DataList[0] = &msg.DMsgData{
		Meta: &meta.DataMeta{
			DataClass: meta.DATA_CLASS_INT32,
			Dimen:     4,
			SampleLen: 1,
		},
		Data: make([]byte, 4*4),
	}

	i := 0
	for i < 4 {
		binary.BigEndian.PutUint32(dmsg.DataList[0].Data[i*4:(i+1)*4], uint32(i))
		i++
	}

	dao.Insert(dmsg)

	sql := fmt.Sprintf("insert into dp_0_0 using dp tags(0,0,0) values(?, 1,2,3,4)")
	dao.TaosCli.Exec(sql, dmsg.Ts)
}

func _task_dao_init(dao *dao.DpDao) {
	dao.Open()
	dao.InitTable(&model.DPoint{
		DataMeta: &meta.DataMeta{
			DataClass: meta.DATA_CLASS_INT32,
			Dimen:     4,
		},
	})

	go _task_dao_query(dao)
	// go _task_insert(dao)
}
