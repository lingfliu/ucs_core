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

	DATA_DIMEN     = 1
	DATA_BYTELEN   = 4
	DATA_SAMPLELEN = 1
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
	tic := time.NewTicker(5 * time.Second) // 每 5 秒生成一次数据
	defer tic.Stop()

	for {
		select {
		case <-sigRun.Done():
			ulog.Log().I("mock", "task is canceled, stopping publish")
			dnode.Cli.Stop()
			return
		case <-tic.C:
			if sigRun.Err() != nil {
				ulog.Log().I("mock", "task is canceled, stopping publish")
				return
			}

			dmsg := &msg.DMsg{
				DNodeId:    dnode.Id,               // 节点ID
				DNodeAddr:  "127.0.0.1:10021",      // 节点地址
				DNodeClass: "tehu_tsi_001",         // 节点类型
				DNodeName:  "DN20",                 // 节点名称
				Ts:         time.Now().UnixMilli(), // 当前时间戳（毫秒）
				Idx:        0,                      // 序号
				Session:    "",                     // 会话标识
				Mode:       0,                      // 模式
				Sps:        20,                     // 采样频率
				SampleLen:  DATA_SAMPLELEN,         // 采样长度
				DataList: []*msg.DMsgData{
					{
						Offset:  0,
						PtAlias: "温度",
						Meta: &meta.DataMeta{
							Dimen:     DATA_DIMEN,
							ByteLen:   4,
							ValAlias:  []string{"Temperature"},
							DataClass: meta.DATA_CLASS_FLOAT,
							Unit:      "°C",
							Msb:       true,
						},
						Data: make([]byte, 4),
					},
					{
						Offset:  1,
						PtAlias: "湿度",
						Meta: &meta.DataMeta{
							Dimen:     DATA_DIMEN,
							ByteLen:   2,
							ValAlias:  []string{"Humidity"},
							DataClass: meta.DATA_CLASS_INT16,
							Unit:      "%",
							Msb:       true,
						},
						Data: make([]byte, 2),
					},
				},
			}

			// 随机生成温度和湿度数据
			tempVal := utils.RandFloat32(0.0, 35.0)
			humiVal := utils.RandInt32(20, 80)
			binary.BigEndian.PutUint32(dmsg.DataList[0].Data, uint32(tempVal))
			binary.BigEndian.PutUint16(dmsg.DataList[1].Data, uint16(humiVal))
			ulog.Log().I("mock", fmt.Sprintf("Generated temp: %.2f°C, humi: %d%%", tempVal, humiVal))

			// 将数据转换为 JSON 并发布到 MQTT
			bytes, err := json.Marshal(dmsg)
			if err != nil {
				ulog.Log().E("mock", "failed to marshal dmsg, err: "+err.Error())
				continue
			}
			ulog.Log().I("mock", "------------------> publish tehu msg: "+string(bytes))
			dnode.Cli.Publish(MQTT_TOPIC, bytes)
		}
	}
}
