package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/lingfliu/ucs_core/dd"
	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

// const mqtt_host = "62.234.16.239:1883"
const mqtt_host = "127.0.0.1:1883"
const mqtt_username = "admin"
const mqtt_password = "admin1234"
const topic = "ucs/dd/mock"

func main() {
	ulog.Config(ulog.LOG_LEVEL_INFO, "", false)

	mqttCli := dd.NewMqttCli(mqtt_host, mqtt_username, mqtt_password)

	sigRun, cancelRun := context.WithCancel(context.Background())
	io := mqttCli.Start()
	go _task_connect(sigRun, io, mqttCli)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	for {
		select {
		case <-s:
			cancelRun()
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func _task_publish(sigRun context.Context, cli *dd.MqttCli) {
	cnt := 0
	for {
		select {
		case <-sigRun.Done():
			return
		default:
			cli.Publish(topic, []byte("hello "+strconv.FormatInt(utils.CurrentTime(), 10)))
			cnt++
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func _task_subscribe(sigRun context.Context, cli *dd.MqttCli) {
	cli.Subscribe(topic)
	for {
		select {
		case <-sigRun.Done():
			return
		case msg := <-cli.RxMsg:
			// println("topic: " + msg.Topic + " data: " + string(msg.Data))
			strs := strings.Split(string(msg.Data), " ")

			tic := strs[len(strs)-1]
			ticInt, err := strconv.ParseInt(tic, 10, 64)
			if err != nil {
				log.Print(err)
			}
			latency := utils.CurrentTime() - ticInt
			println("topic: " + msg.Topic + " pub tic: " + tic + " latency: " + strconv.FormatInt(latency, 10))
		}
	}
}

func _task_connect(sigRun context.Context, io chan int, cli *dd.MqttCli) {
	for {
		select {
		case <-sigRun.Done():
			return
		case state := <-cli.Io:
			if state == dd.DD_STATE_CONNECTED {
				go _task_publish(sigRun, cli)
				go _task_subscribe(sigRun, cli)
			} else if state == dd.DD_STATE_DISCONNECTED {
				cli.Stop()
			}

		}
	}
}
