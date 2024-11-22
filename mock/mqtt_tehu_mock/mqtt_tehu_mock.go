package main

import (
	// "database/sql"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/lingfliu/ucs_core/dao"
	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/dd"
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/model/msg"

	//"github.com/lingfliu/ucs_core/model/spec"
	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
	_ "github.com/taosdata/driver-go/v3/taosSql"
)

type MqttCfg struct {
	Host      string
	Port      int
	Username  string
	Password  string
	TopicList []string
	Qos       byte
	Timeout   int
}

var mqttCfg = &MqttCfg{
	Host:      "62.234.16.239",
	Port:      1883,
	Username:  "admin",
	Password:  "admin1234",
	TopicList: []string{"ucs/dd/tehu_node"},
	Qos:       0,
	Timeout:   3000,
}

const (
	TAOS_HOST     = "62.234.16.239:6030"
	TAOS_DATABASE = "ucs"
	TAOS_USERNAME = "root"
	TAOS_PASSWORD = "taosdata"
)

func main() {
	//config log
	ulog.Config(ulog.LOG_LEVEL_INFO, "./log.log", false)

	// open mqttCli
	mqttCli := dd.NewMqttCli(utils.IpPortJoin(mqttCfg.Host, mqttCfg.Port), mqttCfg.Username, mqttCfg.Password, mqttCfg.TopicList, mqttCfg.Qos, mqttCfg.Timeout)
	mqttCli.Start()
	fmt.Println("MQTT client started")

	//config taos
	dpDataDao := dao.NewDpDataDao(TAOS_HOST, TAOS_DATABASE, TAOS_USERNAME, TAOS_PASSWORD)
	fmt.Println("Opening database connection...")
	go _task_dao_init(dpDataDao)

	sigRun, cancelRun := context.WithCancel(context.Background())
	go _task_mqtt_io(sigRun, mqttCli)
	go _task_recv_mqtt(sigRun, mqttCli, dpDataDao)
	go _task_dao_query(dpDataDao)
	go _task_dao_AggrQuery(dpDataDao)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	for {
		select {
		case <-s:
			cancelRun()
			mqttCli.Stop()
			dpDataDao.Close()
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
func _task_dao_query(dao *dao.DpDataDao) {
	ulog.Log().I("main", "Starting _task_dao_query")
	tic := "2024-11-01 00:00:00.000"
	toc := "2024-11-30 00:00:00.000"

	ulog.Log().I("main", fmt.Sprintf("Querying data from %s to %s", tic, toc))
	tsList, dpData := dao.Query(tic, toc, "tehu_tsi_001", 0, 0, 100, 0, &meta.DataMeta{
		Dimen:   1,
		ByteLen: 4,
	})

	if len(tsList) == 0 || len(dpData.Data) == 0 {
		ulog.Log().I("main", "No data returned from query")
		return
	}

	for i, ts := range tsList {
		val := dpData.Data[i].([]any)[0].(float32)
		ulog.Log().I("main", fmt.Sprintf("ts: %d, %f", ts, val))
	}

}

func _task_dao_AggrQuery(dao *dao.DpDataDao) {
	ulog.Log().I("main", "Starting _task_dao_AggrQuery")

	tic := "2024-11-01 00:00:00.000"
	toc := "2024-11-30 00:00:00.000"
	windowMinutes := int64(10)          // 时间窗口，单位：分钟
	stepMinutes := int64(5)             // 步长，单位：分钟
	window := windowMinutes * 60 * 1000 // 转换为毫秒
	step := stepMinutes * 60 * 1000     // 转换为毫秒
	ops := []int{
		rtdb.TAOS_AGGR_AVG, // 平均值
		rtdb.TAOS_AGGR_MAX, // 最大值
	}

	ulog.Log().I("main", fmt.Sprintf("AggrQuery data between %s and %s interval(%d) sliding(%d)", tic, toc, window, step))
	tsList, dpData := dao.AggrQuery(tic, toc, "tehu_tsi_001", 0, 0, 0, &meta.DataMeta{
		Dimen:    1,
		ByteLen:  4,
		ValAlias: []string{"Temperature"},
	}, window, step, ops)

	if len(tsList) == 0 || len(dpData.Data) == 0 {
		ulog.Log().I("main", "No data returned from aggregate query")
		return
	}

	for i, ts := range tsList {
		ulog.Log().I("main", fmt.Sprintf("dpData.Data[%d]: %v", i, dpData.Data[i]))
		vals := dpData.Data[i].([]float32)
		avg := vals[0] // 平均值
		max := vals[1] // 最大值

		ulog.Log().I("main", fmt.Sprintf("ts: %d, avg: %f, max: %f", ts, avg, max))
	}
}

// 处理 MQTT 客户端的连接状态
func _task_mqtt_io(sigRun context.Context, mqttCli *dd.MqttCli) {
	for state := range mqttCli.Io {
		if state == dd.DD_STATE_CONNECTED {
			ulog.Log().I("main", "mqtt connected")
		} else if state == dd.DD_STATE_DISCONNECTED {
			ulog.Log().I("main", "mqtt disconnected")
		}
	}
}

func _task_dao_init(dao *dao.DpDataDao) {

	template := &model.DNodeTemplate{
		Id:   1,
		Name: "tehu_tsi_001",
		Template: &model.DNode{
			Id:         0,
			TemplateId: 0,
			ParentId:   0,
			Addr:       "0.0.0.0:8008",
			Class:      "tehu_tsi_001",
			Name:       "N20",
			Descrip:    "Tehu sensor",
			PropSet:    make(map[string]string),
			DPointList: []*model.DPoint{
				&model.DPoint{
					Offset: 0,
					Alias:  "Temperature",
					Class:  "温度",
					DataMeta: &meta.DataMeta{
						ByteLen:   4,
						Dimen:     1,
						DataClass: meta.DATA_CLASS_FLOAT,
						Unit:      "°C",
						ValAlias:  []string{"Temperature"},
						Msb:       true,
					},
					Data: make([]byte, 4),
				},
				&model.DPoint{
					Offset: 1,
					Alias:  "Humidity",
					Class:  "湿度",
					DataMeta: &meta.DataMeta{
						ByteLen:   2,
						Dimen:     1,
						DataClass: meta.DATA_CLASS_INT16,
						Unit:      "%",
						ValAlias:  []string{"Humidity"},
						Msb:       true,
					},
					Data: make([]byte, 2),
				},
			},
			State:     0,
			Sps:       20,
			SampleLen: 2,
			Mode:      model.DNODE_MODE_AUTO,
		},
	}
	dao.Open()

	res := dao.CreateTableFromTemplate(template)
	if res < 0 {
		ulog.Log().E("mock_taos", "Create stable failed")
	} else {
		ulog.Log().I("mock_taos", "create stable success")
	}
}

func _task_recv_mqtt(sigRun context.Context, mqttCli *dd.MqttCli, dpDataDao *dao.DpDataDao) {
	for _, topic := range mqttCfg.TopicList {
		if result := mqttCli.Subscribe(topic); result != 0 {
			log.Println("Failed to subscribe to topic:", topic)
		} else {
			log.Println("Successfully subscribed to topic:", topic)
		}
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigRun.Done():
			fmt.Println("Signal received, terminating loop.")
			return
		case rxmsg := <-mqttCli.RxMsg:
			fmt.Println("=================接收MQTT消息=================")
			fmt.Printf("Received message on topic: " + rxmsg.Topic + ", Data: " + string(rxmsg.Data))
			switch rxmsg.Topic {
			case "ucs/dd/tehu_node":
				dpMsg := &msg.DMsg{}
				err := json.Unmarshal(rxmsg.Data, dpMsg)
				if err != nil {
					ulog.Log().E("main", "tehu_node msg decode error: "+err.Error())
				} else {
					ulog.Log().I("main", fmt.Sprintf("Decoded dpMsg: %+v", dpMsg))
					ulog.Log().I("main", "insert tehu_node msg: "+string(rxmsg.Data))
					dpDataDao.Insert(dpMsg)
				}
			default:
				log.Printf("Message received on unexpected topic: %s, Data: %s", rxmsg.Topic, string(rxmsg.Data))
			}
		case <-ticker.C:
			fmt.Println("No message received in the last 5 seconds.")
		default:
			time.Sleep(time.Second * 1)
		}
	}
}
