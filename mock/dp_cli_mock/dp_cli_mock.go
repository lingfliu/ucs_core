package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
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
	MQTT_TOPIC    = "ucs/dd/dp"
	MQTT_USERNAME = "admin"
	MQTT_PASSWORD = "admin1234"

	DATA_DIMEN   = 4
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
	for i := 0; i < 10; i++ {
		//定义10个DNode，每个DNode包含2个DPoint
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
	time.Sleep(10 * time.Second)
	for _, ctlRun := range ctlRunList {
		ctlRun()
	}
	ulog.Log().I("mock", "all tasks canceled")
}

func _task_mock_mqtt(sigRun context.Context, dnode *DNodeMock) {
	dnode.Cli.Start()
	tic := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-sigRun.Done():
			dnode.Cli.Stop()
		case <-tic.C:
			dmsg := &msg.DMsg{
				Ts:      time.Now().UnixNano() / 1e6, //in milliseconds
				Sps:     100 * 1000 * 1000,
				Mode:    0,
				DNodeId: dnode.Id,
				//random int32
				DataSet: make(map[int]*msg.DMsgData),
			}

			for i := 0; i < dnode.NoDp; i++ {
				//模拟部分：sampleLen=1，dimen=4，bytelen=4， MSB，随机数据
				valueList := make([]byte, dnode.DpDimen*4)
				for i := 0; i < dnode.DpDimen; i++ {
					//random int32
					v := utils.RandInt32(0, 10000)
					binary.BigEndian.PutUint32(valueList[i*4:(i+1)*4], uint32(v))
				}
				dmsg.DataSet[i] = &msg.DMsgData{
					Meta: &meta.DataMeta{
						Dimen:     DATA_DIMEN,
						SampleLen: 1,
						ByteLen:   DATA_BYTELEN,
						DataClass: meta.DATA_CLASS_INT32,
					},
					Data: valueList,
				}
			}
			bytes, err := json.Marshal(dmsg)
			if err != nil {
				ulog.Log().E("mock", "failed to marshal ddMsg, err: "+err.Error())
			}
			ulog.Log().I("mock", "publish dp msg: "+string(bytes))
			dnode.Cli.Publish(MQTT_TOPIC, bytes)
		}
	}
}
