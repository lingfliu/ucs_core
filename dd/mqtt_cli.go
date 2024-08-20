package dd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lingfliu/ucs_core/utils"
)

type MqttCli struct {
	BaseDdCli

	mc mqtt.Client
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

	return &MqttCli{
		BaseDdCli: BaseDdCli{
			Host: host,
		},
		mc: mc,
	}
}

func (cli *MqttCli) Subscribe(topic string) {
}

func (cli *MqttCli) Unsubscribe(topic string) {
}

func (cli *MqttCli) Publish(topic string, data []byte) {
}

func (cli *MqttCli) Connect() {
	token := cli.mc.Connect()
	if token.WaitTimeout(time.Second*3) && token.Error() != nil {
		panic(token.Error())
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
