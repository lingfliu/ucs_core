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
	"github.com/lingfliu/ucs_core/dao/impl"
	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/dd"
	"github.com/lingfliu/ucs_core/model"
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
	taosCli := rtdb.NewTaosCli(TAOS_HOST, TAOS_DATABASE, TAOS_USERNAME, TAOS_PASSWORD)
	log.Println("Opening database connection...")
	taosCli.Open()

	// 实例化 TehuNodeDao
	template := model.DPoint{}
	colNameList := []string{"temp", "humi"}
	dpDao := &dao.DpDao{
		TaosCli:     taosCli,
		Template:    &template,
		ColNameList: colNameList,
	}
	tehuNodeDao := &impl.TehuNodeDao{
		DpDao: *dpDao,
	}
	fmt.Printf("TehuNodeDao created: %+v\n", tehuNodeDao)

	node := tehuNodeDao.GenerateTemplate()
	fmt.Printf("Generated node: %+v\n", node)

	// 调用 Create 函数
	res := tehuNodeDao.Create()
	if res < 0 {
		log.Println("Failed to create stable.")
	} else {
		log.Println("Stable created successfully.")
		defer dpDao.TaosCli.Close() // 确保在退出时关闭连接
	}
	//调用 TableExist 方法
	exists := tehuNodeDao.TableExist()
	if exists {
		fmt.Println("所有表存在")
	} else {
		fmt.Println("有表不存在")
	}

	sigRun, cancelRun := context.WithCancel(context.Background())
	go _task_mqtt_io(sigRun, mqttCli)

	fmt.Println("========= MQTT 客户端连接成功===============")

	go _task_recv_mqtt(sigRun, mqttCli, tehuNodeDao)

	//等待退出信号，并在收到信号时断开连接
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

func _task_recv_mqtt(sigRun context.Context, mqttCli *dd.MqttCli, tehuNodeDao *impl.TehuNodeDao) {
	//Subscribe( mqtt message
	for _, topic := range mqttCfg.TopicList {
		if result := mqttCli.Subscribe(topic); result != 0 {
			log.Println("Failed to subscribe to topic:", topic)
		} else {
			log.Println("Successfully subscribed to topic:", topic)
		}
	}
	// Add a ticker for logging if no messages are received
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigRun.Done():
			fmt.Println("Signal received, terminating loop.")
			return
		case rxmsg := <-mqttCli.RxMsg:
			fmt.Println("=========接收MQTT消息=================")
			fmt.Printf("Received message on topic: " + rxmsg.Topic + ", Data: " + string(rxmsg.Data))
			//ulog.Log().I("main", "Received message on topic: "+rxmsg.Topic+", Data: "+string(rxmsg.Data))
			//parse mqtt message
			switch rxmsg.Topic {
			case "ucs/dd/tehu_node":
				//payload is encoded by DpCoder: [ts, dnode id, dp offsetidx, int value]
				dpMsg := &msg.DMsg{}
				err := json.Unmarshal(rxmsg.Data, dpMsg)
				if err != nil {
					ulog.Log().E("main", "th_node msg decode error: "+err.Error())
				} else {
					ulog.Log().I("main", fmt.Sprintf("Decoded dpMsg: %+v", dpMsg))
					ulog.Log().I("main", fmt.Sprintf("DNodeId: %d, Offset: %d, DataSet: %+v", dpMsg.DNodeId, dpMsg.Offset, dpMsg.DataSet))
					//insert into taos
					ulog.Log().I("main", "insert tehu_node msg: "+string(rxmsg.Data))
					tehuNodeDao.Insert(dpMsg)
					//dpDao.Insert(dpMsg)
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
