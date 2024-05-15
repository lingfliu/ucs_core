package main

import (
	"sync"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
)

var wg sync.WaitGroup

func main() {
	//config log
	ulog.Config(ulog.LOG_LEVEL_INFO, "./log.log", false)
	// //service initialization
	// //start servers
	// agvSrv := srv.CreateTcpServer(cfg.port, cfg.addr)
	// wsnSrv := srv.CreateTcpServer(cfg.port, cfg.addr)
	// nvsSrv := srv.CreateTcpServer(cfg.port, cfg.addr)

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

	//wait for close signal
	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	rx := make(chan int, 10)

	wg.Add(1)
	go _task_consume(rx)
	wg.Add(1)
	go _task_feed(rx)

	wg.Wait()

}

func _task_consume(rx chan int) {
	defer wg.Done()
	for {
		b := <-rx
		// print(b)
		output := struct {
			Value int
		}{
			Value: b,
		}

		ulog.Log().W("main", output)
	}
}
func _task_feed(rx chan int) {
	defer wg.Done()

	cnt := 0
	tic := time.NewTicker(time.Millisecond * 10)
	for range tic.C {
		rx <- cnt
		cnt++
	}
}
