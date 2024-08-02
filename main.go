package main

import (
	"context"
	"sync"

	"github.com/lingfliu/ucs_core/ulog"
)

var wg sync.WaitGroup

//global variables

func main() {
	//config log
	ulog.Config(ulog.LOG_LEVEL_INFO, "./log.log", false)

	// //service initialization

	// //membuff
	// memBuff := membuff.CreateMemBuff(cfg.memBuffSize)

	// //mq
	// //1. inmem-mq: memq
	// meMq := mq.CreateMq(cfg.mqAddr, cfg.mqPort, cfg.mqUser, cfg.mqPwd, cfg.mqVhost)
	// //2. emqx
	// mqttCli := cli.CreateMqttClient(cfg.mqttAddr, cfg.mqttPort, cfg.mqttUser, cfg.mqttPwd)
	// //3. nsq mq
	// nsqMq := mq.CreateNsqMq(cfg.nsqAddr, cfg.nsqPort)
	// //4. kafka mq
	// kafkaMq := mq.CreateKafkaMq(cfg.kafkaAddr, cfg.kafkaPort)

	// //dba
	// redisCli := dba.CreateRedisClient(cfg.redisAddr, cfg.redisPort, cfg.redisPwd, cfg.redisDb)
	// sqlCli := dba.CreateSqlClient(cfg.sqlAddr, cfg.sqlPort, cfg.sqlUser, cfg.sqlPwd, cfg.sqlDb)
	// mongoCli = dba.CreateMongoClient(cfg.mongoAddr, cfg.mongoPort, cfg.mongoUser, cfg.mongoPwd, cfg.mongoDb)
	// minioCli = dba.CreateMinioClient(cfg.minioAddr, cfg.minioPort, cfg.minioUser, cfg.minioPwd)

	//service discovery
	//start go-zero services

	// global context
	// ctxGlobal := context.Background()

	// start conn servers
	// srvMach := spec.NewMachServer(cfg.mach.port)
	// dataCh := srvMach.Expose()
	// go _task_mach_data_process(dataCh)
	// srvMach.Start()
	// srvIot := spec.NewIotServer(cfg.mqtt)
	// srvIot.Start()
	// srvStream := spec.NewNsServer(cfg.port)
	// srvStream.Start()

	// //start cloud-edge sync service
	// serviceCesync := spec.NewCesyncService(cfg.host, cfg.host_ports)
}

func _task_mach_data_process(ctx context.Context, ch chan MachData) {
	for {
		select {
		case <-ch:
			go func() {
				// serviceMachDataProcess.
				// 	Filter(data).
				// 	Then().
				// 	Align().
				// 	Then().
				// 	Publish()
			}()
		case <-ctx.Done():
			break
		}
	}
}

// func _task_consume(rx chan int) {
// 	defer wg.Done()
// 	for {
// 		b := <-rx
// 		// print(b)
// 		output := struct {
// 			Value int
// 		}{
// 			Value: b,
// 		}

// 		ulog.Log().W("main", output)
// 	}
// }
// func _task_feed(rx chan int) {
// 	defer wg.Done()

// 	cnt := 0
// 	tic := time.NewTicker(time.Millisecond * 10)
// 	for range tic.C {
// 		rx <- cnt
// 		cnt++
// 	}
// }
