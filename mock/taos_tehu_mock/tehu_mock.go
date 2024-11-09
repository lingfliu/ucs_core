package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/lingfliu/ucs_core/dd"
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/model/msg"
	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

const (
	MQTT_HOST     = "62.234.16.239"
	MQTT_PORT     = 1883
	MQTT_TOPIC    = "ucs/dd/tehu_node"
	MQTT_USERNAME = "admin"
	MQTT_PASSWORD = "admin1234"

	DATA_DIMEN   = 1
	DATA_BYTELEN = 4 //int16
)

type DNodeMock struct {
	Class   int64
	Id      int64
	NoDp    int
	DpDimen int
	Cli     *dd.MqttCli
}

func main() {
	//config log
	ulog.Config(ulog.LOG_LEVEL_INFO, "./log.log", false)

	nodeList := make([]*DNodeMock, 0)
	ctlRunList := make([]context.CancelFunc, 0)
	for i := 0; i < 5; i++ {
		//定义5个DNode，每个DNode包含2个DPoint
		cli := dd.NewMqttCli(utils.IpPortJoin(MQTT_HOST, MQTT_PORT),
			MQTT_USERNAME,
			MQTT_PASSWORD,
			[]string{MQTT_TOPIC},
			0,
			3000)
		dnode := &DNodeMock{
			Class:   int64(i),
			Id:      int64(i),
			NoDp:    2,
			DpDimen: DATA_DIMEN,
			Cli:     cli,
		}

		nodeList = append(nodeList, dnode)
		sigRun, ctlRun := context.WithCancel(context.Background())
		ctlRunList = append(ctlRunList, ctlRun)
		go _task_mock_mqtt(sigRun, dnode)
	}

	go _task_cancel(ctlRunList)

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

func _task_cancel(ctlRunList []context.CancelFunc) {
	time.Sleep(15 * time.Second)
	for _, ctlRun := range ctlRunList {
		ctlRun()
	}
	ulog.Log().I("mock", "all tasks canceled")
}

func _task_mock_mqtt(sigRun context.Context, dnode *DNodeMock) {
	dnode.Cli.Start()
	tic := time.NewTicker(5 * time.Second) //每5秒生成一次数据

	// go func() {
	// 	dnode.Cli.Subscribe(MQTT_TOPIC)
	// }()

	for {
		select {
		case <-sigRun.Done():
			dnode.Cli.Stop()
		case <-tic.C:
			if sigRun.Err() != nil {
				ulog.Log().I("mock", "task is canceled, stopping publish")
				return
			}
			for i := 0; i < dnode.NoDp; i++ {
				//模拟部分：NoDp=2, sampleLen=1，dimen=1，bytelen=4， MSB，随机数据
				dmsg := &msg.DMsg{
					Ts:      time.Now().UnixNano() / 1e6, //当前时间戳，单位为毫秒
					Sps:     100 * 1000 * 1000,           //采样率
					Mode:    0,                           //0-定时采样，1-事件触发，2-轮询
					DNodeId: dnode.Id,                    // 节点ID
					DataSet: make(map[int]*msg.DMsgData),
				}

				valueList := make([]byte, dnode.DpDimen*4)
				if i == 0 {
					// 模拟生成温度数据 (维度0)
					temp := utils.RandInt32(0, 50)                              // 随机生成温度范围0°C到50°C
					binary.BigEndian.PutUint32(valueList[0:4], uint32(temp))    // 将温度值放入第0个维度
					ulog.Log().I("mock", fmt.Sprintf("publish temp: %d", temp)) // 打印温度值

					dmsg.Offset = 0
					dmsg.DataSet[dmsg.Offset] = &msg.DMsgData{
						Meta: &meta.DataMeta{
							Dimen:     1, // 维度
							SampleLen: 1, // 采样长度
							ByteLen:   4, // uint32 占用 4 字节
							DataClass: meta.DATA_CLASS_INT,
							Alias:     "Temperature",
							Unit:      "°C",
						},
						Data: valueList,
					}

				} else if i == 1 {
					// 模拟生成湿度数据 (维度1)
					humi := utils.RandInt32(20, 70)                             // 随机生成湿度范围20%到70%
					binary.BigEndian.PutUint32(valueList[0:4], uint32(humi))    // 将湿度值放入第1个维度
					ulog.Log().I("mock", fmt.Sprintf("publish humi: %d", humi)) // 打印湿度值

					dmsg.Offset = 1
					dmsg.DataSet[dmsg.Offset] = &msg.DMsgData{
						Meta: &meta.DataMeta{
							Dimen:     DATA_DIMEN, //1
							SampleLen: 1,
							ByteLen:   DATA_BYTELEN, // uint32 占用 4 字节
							DataClass: meta.DATA_CLASS_INT,
							Alias:     "Humidity",
							Unit:      "%",
						},
						Data: valueList,
					}
					//指定维度，不用循环
					// for i := 0; i < dnode.DpDimen; i++ {
					// 	//random int32
					// 	v := utils.RandInt32(0, 10000)
					// 	binary.BigEndian.PutUint32(valueList[i*4:(i+1)*4], uint32(v))
					// 	// 打印温湿度值
					// 	ulog.Log().I("mock", fmt.Sprintf("publish temp and humi: %d", uint32(v)))
					// }

				}
				// 将数据转换为 JSON 并发布到 MQTT
				bytes, err := json.Marshal(dmsg)
				if err != nil {
					ulog.Log().E("mock", "failed to marshal ddMsg, err: "+err.Error())
					continue
				}
				ulog.Log().I("mock", "------------------>publish tehu msg: "+string(bytes))
				dnode.Cli.Publish(MQTT_TOPIC, bytes)
				//fmt.Println("Publishing to MQTT topic: ucs/dd/th_node, Data:", bytes)

			}
		}
	}
}
