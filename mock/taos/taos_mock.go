package main

import (
	"encoding/binary"
	"os"
	"os/signal"
	"time"

	"github.com/lingfliu/ucs_core/dao"
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/model/msg"
	"github.com/lingfliu/ucs_core/ulog"
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
	dpDataDao := dao.NewDpDataDao(TAOS_HOST, TAOS_DATABASE, TAOS_USERNAME, TAOS_PASSWORD)
	go _task_dao_init(dpDataDao)
	// go _task_insert(dpDataDao)

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

func _task_insert(dao *dao.DpDataDao) {
	idx := 0
	tic := time.Tick(20 * time.Millisecond)
	for {

		select {
		case <-tic:
			dmsg := &msg.DMsg{
				DNodeId:    0, //should provide at least id / name / addr
				DNodeAddr:  "127.0.0.1:10021",
				DNodeClass: "tehu_tsi_002",         //if class is missing, receiver should look for the class from db
				DNodeName:  "DN21",                 //if name is missing, receiver should look for the name from db
				Ts:         time.Now().UnixMilli(), //timestamp of first sample
				Idx:        idx,                    //序号， 用于辅助判断是否丢包
				Session:    "",                     //会话标识
				Mode:       0,                      //模式，对应DNode的Mode: 0-定时采样，1-事件触发，2-轮询
				Sps:        20,                     //采样频率, 仅在Mode=0时有效
				SampleLen:  1,                      //采样长度
				DataList: []*msg.DMsgData{
					&msg.DMsgData{
						Offset:  0,
						PtAlias: "温度",
						Meta: &meta.DataMeta{
							Dimen:     1,
							ByteLen:   4,
							DataClass: meta.DATA_CLASS_FLOAT,
							Unit:      "C",
							Msb:       true,
						},
						Data: make([]byte, 4),
					},
					&msg.DMsgData{
						Offset:  1,
						PtAlias: "湿度",
						Meta: &meta.DataMeta{
							Dimen:     1,
							ByteLen:   2,
							DataClass: meta.DATA_CLASS_INT16,
							Unit:      "H",
							Msb:       true,
						},
						Data: make([]byte, 2),
					},
				}, //offset as the key
			}
			tempVal := 25.0
			humiVal := 70
			binary.BigEndian.PutUint32(dmsg.DataList[0].Data, uint32(tempVal))
			binary.BigEndian.PutUint16(dmsg.DataList[1].Data, uint16(humiVal))
			dao.Insert(dmsg)
			idx++
		default:
			time.Sleep(1 * time.Second)
		}
	}

}

func _task_dao_init(dao *dao.DpDataDao) {

	template := &model.DNodeTemplate{
		Id:   1,
		Name: "tehu_tsi_002",
		Template: &model.DNode{
			Id:         0,
			TemplateId: 0,
			ParentId:   0,
			Addr:       "0.0.0.0:8008",
			Class:      "tehu_tsi_002",
			Name:       "N21",
			Descrip:    "Tehu sensor",
			PropSet:    make(map[string]string),
			DPointList: []*model.DPoint{
				&model.DPoint{
					Offset: 0,
					Alias:  "",
					Class:  "温度",
					DataMeta: &meta.DataMeta{
						ByteLen:   4,
						Dimen:     1,
						DataClass: meta.DATA_CLASS_FLOAT,
						Unit:      "C",
					},
					Data: make([]byte, 4),
				},
				&model.DPoint{
					Offset: 1,
					Alias:  "",
					Class:  "湿度",
					DataMeta: &meta.DataMeta{
						ByteLen:   2,
						Dimen:     1,
						DataClass: meta.DATA_CLASS_INT16,
						Unit:      "%",
					},
					Data: make([]byte, 2),
				},
			},
			State:     0,
			Sps:       1000,
			SampleLen: 1,
			Mode:      model.DNODE_MODE_AUTO,
		},
	}
	dao.Open()

	res := dao.CreateTableFromTemplate(template)
	if res < 0 {
		ulog.Log().E("mock_taos", "Create table failed")
	} else {
		ulog.Log().I("mock_taos", "create table success")
	}

	// go _task_dao_query(dao)
	go _task_insert(dao)
}
