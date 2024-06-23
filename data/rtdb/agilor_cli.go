package rtdb

const (
	AGILOR_STATE_DISCONNECTED = 0
	AGILOR_STATE_CONNECTED    = 1
)

type AgilorSubscription struct {
	msg   <-chan string
	topic string
}
type AgilorCli struct {
	SrvAddr  string
	SrvPort  int
	Username string
	Password string
	State    int

	SubscriptionList map[string]*AgilorSubscription
}

func NewAgilorCli(srvAddr string, srvPort int, username string, password string) *AgilorCli {
	return &AgilorCli{
		SrvAddr:  srvAddr,
		SrvPort:  srvPort,
		Username: username,
		Password: password,
		State:    AGILOR_STATE_DISCONNECTED,
	}
}

func (cli *AgilorCli) Connect() {
	//TODO: connect using C API
}

func (cli *AgilorCli) Disconnect() {
}

func (cli *AgilorCli) Update() {
}

func (cli *AgilorCli) Subscribe(topic string) *AgilorSubscription {
	msg := make(chan string)
	cli.SubscriptionList[topic] = &AgilorSubscription{
		msg:   msg,
		topic: topic,
	}
	return cli.SubscriptionList[topic]
}

func (cli *AgilorCli) Unsubscribe(topic string) {
	delete(cli.SubscriptionList, topic)
}
