package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/dao"
	"github.com/lingfliu/ucs_core/dd"
	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

// global variables
type MqttCfg struct {
	Host      string
	Port      int
	Username  string
	Password  string
	TopicList []string
	Qos       byte
	Timeout   int
}

type TaosCfg struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

var mqttCfg = &MqttCfg{
	Host:      "62.234.16.239",
	Port:      1883,
	Username:  "admin",
	Password:  "admin1234",
	TopicList: []string{"ucs/dd/dp"},
	Qos:       0,
	Timeout:   3000,
}

var taosCfg = &TaosCfg{
	Host:     "62.234.16.239",
	Port:     6030,
	Username: "root",
	Password: "taosdata",
	DbName:   "ucs",
}

func main() {
	//config log
	ulog.Config(ulog.LOG_LEVEL_INFO, "./log.log", false)

	// open mqttCli
	mqttCli := dd.NewMqttCli(utils.IpPortJoin(mqttCfg.Host, mqttCfg.Port), mqttCfg.Username, mqttCfg.Password, mqttCfg.TopicList, mqttCfg.Qos, mqttCfg.Timeout)
	mqttCli.Start()
	dpDao := dao.NewDpDao(utils.IpPortJoin(taosCfg.Host, taosCfg.Port), taosCfg.DbName, taosCfg.Username, taosCfg.Password)

	go _task_dao_init(dpDao)

	sigRun, cancelRun := context.WithCancel(context.Background())
	go _task_mqtt_io(sigRun, mqttCli)
	go _task_recv_mqtt(sigRun, mqttCli, dpDao)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	for {
		select {
		case <-s:
			cancelRun()
			mqttCli.Stop()
			dpDao.Close()
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}

}

func _task_mqtt_io(sigRun context.Context, mqttCli *dd.MqttCli) {
	for {
		select {
		case state := <-mqttCli.Io:
			if state == dd.DD_STATE_CONNECTED {
				ulog.Log().I("main", "mqtt connected")
			} else if state == dd.DD_STATE_DISCONNECTED {
				ulog.Log().I("main", "mqtt disconnected")
			}
		}
	}
}

func _task_dao_init(dao *dao.DpDao) {
	dao.Open()
	dao.Init()
}

func _task_recv_mqtt(sigRun context.Context, mqttCli *dd.MqttCli, dpDao *dao.DpDao) {
	//task: receive mqtt message
	for {
		select {
		case <-sigRun.Done():
			return
		case msg := <-mqttCli.RxMsg:
			//parse mqtt message
			switch msg.Topic {
			case "ucs/dd/dp":
				//payload is encoded by DpCoder: [ts, dnode id, dp offsetidx, int value]
				dpMsg := &coder.DpMsg{}
				err := json.Unmarshal(msg.Data, dpMsg)
				if err != nil {
					ulog.Log().E("main", "dp msg decode error: "+err.Error())
				} else {
					//insert into taos
					ulog.Log().I("main", "insert dp msg: "+string(msg.Data))
					dpDao.Insert(dpMsg)
				}

			}
		}
	}
}
