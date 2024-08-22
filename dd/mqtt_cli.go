package dd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lingfliu/ucs_core/ulog"
	"github.com/lingfliu/ucs_core/utils"
)

type MqttCli struct {
	BaseDdCli

	Timeout int //timeout in ms
	State   int
	Qos     byte
	mc      mqtt.Client
	opt     *mqtt.ClientOptions

	sigRun    context.Context
	cancelRun context.CancelFunc
}

func NewMqttCli(host string, username string, password string, topicList []string, qos byte, timeout int) *MqttCli {
	cliId := utils.GenMqttCliId()
	opt := mqtt.NewClientOptions()
	opt.AddBroker(host)
	opt.SetUsername(username)
	opt.SetPassword(password)
	opt.SetClientID(cliId)
	opt.SetKeepAlive(time.Second * 60)
	opt.SetAutoReconnect(true)
	opt.SetCleanSession(false) //autoconnect and keep the session
	// opt.SetTLSConfig(loadTLSConfig("caFilePath"))

	sigRun, cancelRun := context.WithCancel(context.Background())
	return &MqttCli{
		BaseDdCli: BaseDdCli{
			Host:      host,
			TopicList: topicList,
			State:     DD_STATE_DISCONNECTED,

			RxMsg: make(chan *DdZeroMsg, 32),
			TxMsg: make(chan *DdZeroMsg, 32),
			Io:    make(chan int),
		},
		Qos: qos,
		mc:  nil,
		opt: opt,

		sigRun:    sigRun,
		cancelRun: cancelRun,
		Timeout:   timeout,
	}
}

func (cli *MqttCli) Connect() {
	if cli.State == DD_STATE_CONNECTED {
		return
	}

	cli.State = DD_STATE_CONNECTING
	cli.mc = mqtt.NewClient(cli.opt)
	token := cli.mc.Connect()
	if token.WaitTimeout(time.Second*3) && token.Error() != nil {
		ulog.Log().E("mqttcli", "fail to connect")
		cli.State = DD_STATE_DISCONNECTED
	} else {
		ulog.Log().I("mqttcli", "connected")
		cli.State = DD_STATE_CONNECTED
		cli.Io <- DD_STATE_CONNECTED
		go cli._task_subscribe()
	}
}

func (cli *MqttCli) Disconnect() {
	cli.State = DD_STATE_DISCONNECTED
	cli.Io <- DD_STATE_DISCONNECTED
	cli.mc.Disconnect(250)
}

func (cli *MqttCli) Subscribe(topic string) int {
	ulog.Log().I("mqtt", "subscribed to "+topic)

	if cli.State != DD_STATE_CONNECTED {
		cli.Connect()
	}

	token := cli.mc.Subscribe(topic, cli.Qos, func(c mqtt.Client, m mqtt.Message) {
		cli.RxMsg <- &DdZeroMsg{
			Topic: m.Topic(),
			Data:  m.Payload(),
		}
	})

	if token.WaitTimeout(time.Duration(2000*time.Millisecond)) && token.Error() != nil {
		cli.Disconnect()
		return -1
	}
	return 0
}

func (cli *MqttCli) Unsubscribe(topic string) int {
	for idx, t := range cli.TopicList {
		if t == topic {
			cli.TopicList = append(cli.TopicList[:idx], cli.TopicList[idx+1:]...)

			if token := cli.mc.Unsubscribe(topic); token.WaitTimeout(time.Duration(cli.Timeout*int(time.Millisecond))) && token.Error() != nil {
				ulog.Log().E("mqttcli", "unsubscribe error: "+token.Error().Error())
				go cli.Disconnect()
				return -1
			}
			return 0
		}
	}
	return -2 //topic not in the list
}

func (cli *MqttCli) Publish(topic string, data []byte) {
	if cli.State != DD_STATE_CONNECTED {
		cli.Connect()
	}

	if token := cli.mc.Publish(topic, cli.Qos, false, data); token.WaitTimeout(time.Duration(cli.Timeout*int(time.Millisecond))) && token.Error() != nil {
		cli.Disconnect()
	}
}

func (cli *MqttCli) Start() {
	go cli._task_connect()
}

func (cli *MqttCli) Stop() {
	cli.cancelRun()
	cli.Disconnect()
}

func (cli *MqttCli) _task_connect() {
	tic := time.NewTicker(time.Millisecond * 200)
	cli.Connect()
	for {
		select {
		case <-cli.sigRun.Done():
			cli.Stop()
		case <-tic.C:
			if cli.State == DD_STATE_DISCONNECTED {
				cli.Connect()
			}
		}
	}
}

func (cli *MqttCli) _task_subscribe() {

	for _, topic := range cli.TopicList {
		ret := cli.Subscribe(topic)
		if ret < 0 {
			return
		}
	}
}

// TODO: add tls support
func loadTLSConfig(caFile string) *tls.Config {
	// load tls config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	if caFile != "" {
		certpool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	return &tlsConfig
}
