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

	State int
	Qos   int
	mc    mqtt.Client
	RxMsg chan *DdMsg
	TxMsg chan *DdMsg
	Io    chan int

	sigRun    context.Context
	cancelRun context.CancelFunc
}

func NewMqttCli(host string, username string, password string) *MqttCli {
	cliId := utils.GenMqttCliId()
	opt := mqtt.NewClientOptions()
	opt.AddBroker(host)
	opt.SetUsername(username)
	opt.SetPassword(password)
	opt.SetClientID(cliId)
	opt.SetKeepAlive(time.Second * 60)
	// opt.SetTLSConfig(loadTLSConfig("caFilePath"))

	mc := mqtt.NewClient(opt)

	sigRun, cancelRun := context.WithCancel(context.Background())
	return &MqttCli{
		BaseDdCli: BaseDdCli{
			Host:     host,
			TopicSet: make(map[string]string),
			State:    DD_STATE_DISCONNECTED,
		},
		RxMsg:     make(chan *DdMsg, 32),
		TxMsg:     make(chan *DdMsg, 32),
		Io:        make(chan int),
		Qos:       0,
		mc:        mc,
		sigRun:    sigRun,
		cancelRun: cancelRun,
	}
}

func (cli *MqttCli) Subscribe(topic string) int {
	if _, ok := cli.TopicSet[topic]; ok {
		return -2
	}

	token := cli.mc.Subscribe(topic, byte(cli.Qos), func(c mqtt.Client, m mqtt.Message) {
		cli.RxMsg <- &DdMsg{
			Topic: m.Topic(),
			Data:  m.Payload(),
		}
	})
	if token.Wait() && token.Error() != nil {
		cli.Disconnect()
		return -1
	}
	cli.TopicSet[topic] = topic
	return 0
}

func (cli *MqttCli) Unsubscribe(topic string) int {
	if token := cli.mc.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		ulog.Log().E("mqttcli", "unsubscribe error: "+token.Error().Error())
		go cli.Disconnect()
		return -1
	}
	return 0
}

func (cli *MqttCli) Publish(topic string, data []byte) {
	cli.mc.Publish(topic, 0, false, data)
}

func (cli *MqttCli) Connect() {
	cli.State = DD_STATE_CONNECTING
	token := cli.mc.Connect()
	if token.WaitTimeout(time.Second*3) && token.Error() != nil {
		// panic(token.Error())
		cli.State = DD_STATE_DISCONNECTED
		cli.Io <- DD_STATE_DISCONNECTED
	} else {
		ulog.Log().I("mqttcli", "connected")
		cli.State = DD_STATE_CONNECTED
		cli.Io <- DD_STATE_CONNECTED
	}

}

func (cli *MqttCli) Disconnect() {
	cli.State = DD_STATE_DISCONNECTED
	cli.Io <- DD_STATE_CONNECTED
	cli.mc.Disconnect(250)
}

func (cli *MqttCli) Start() chan int {
	go cli.Connect()
	return cli.Io
}

func (cli *MqttCli) Stop() {
	cli.cancelRun()
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
