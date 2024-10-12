package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/lingfliu/ucs_core/dao"
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
	dao.Query(tic, toc)
}

func _task_dao_init(dao *dao.DpDao) {
	dao.Open()
	dao.Init()

	go _task_dao_query(dao)
}
