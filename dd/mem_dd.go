package dd

/**
 * int hash function for topic
 */
func topic2int(topic string) uint64 {
	return 0
}

type DdStream struct {
	Topic  int
	Stream chan []byte
}

type DdMsg struct {
	Topic int
	Data  []byte
}

type MemDdCli struct {
	// ConnCli  *conn.Cli
	//connection
	Host     string
	Prop     int
	Username string
	Password string

	TopicSet map[int]bool
}

func (cli *MemDdCli) Connect() {
}

func (cli *MemDdCli) Disconnect() {
}

func (cli *MemDdCli) Subscribe(topic string) {
}

const (
	TRANSPORT_TCP  = "tcp"
	TRANSPORT_UDP  = "udp"
	TRANSPORT_QUIC = "quic"
)

type MemDd struct {
	Qos            int    //0, 1, 2
	TransportProto string //tcp, udp, quic, kcp
	TopicSet       map[int]*string
	PublisherSet   map[string]*MemDdCli
	SubscriberSet  map[string]*MemDdCli
	MsgSet         map[int]chan []byte
	StreamSet      map[int]*DdStream
}

func (dd *MemDd) Listen() {
}

func (dd *MemDd) OnSubscribe(topicId int, cli *MemDdCli) {
}

func (dd *MemDd) OnUnsubscribe(topicId int) {
}

func (dd *MemDd) Subscribe(topicId int, cli *MemDdCli) {
	// if _, ok := dd.SubscriberSet[cli.C.RemoteAddr]; !ok {
	// }
}

func (dd *MemDd) Publish(msg *DdMsg) {
}
