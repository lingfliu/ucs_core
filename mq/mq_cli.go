package mq

import {
	"github.com/nsqio/go-nsq"
}
/**
 * Message queue client for nsq
 */
type MqMsg struct {
	content string
}


type MQCli struct {
	producer *nsq.Producer
	consumer *nsq.Consumer
}

func NewMqCli(mq_addr string) *MQCli {
	config := nsq.NewConfig()
	config.MaxInFlight = 9
	producter, err := nsq.NewProducer(mq_addr, config)
	if err != nil {
		log.Fatal("create nsq producer failed")
	}

	return &MQCli{
		producer: producter,
		consumer: nil,
	}
}

func (cli *MQCli) Shutdown(topic string) {
	cli.producer.Stop()
}

func (cli *MQCli) PublishJson(json string, topic string) {
	cli.producer.Publish(topic, []byte(json))
}

func (cli *MQCli) PublishObj(msg any, topic string) {
	cli.producer.Publish(topic, msg.Serialize())
}

func (cli *MQCli) PublishStream(url string, topic string) {
	cli.producer.Publish(topic, []byte(url))
}
