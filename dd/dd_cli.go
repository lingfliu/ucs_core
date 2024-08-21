package dd

const (
	DD_STATE_DISCONNECTED = 0
	DD_STATE_CONNECTING   = 1
	DD_STATE_CONNECTED    = 2
	DD_STATE_CLOSE        = 3
)

type DdCli interface {
	Connect() int
	Disconnect() int
	Close() int
	Subscribe(topic string) int64
	Unsubscribe(topic string) int64
	Publish(topic string, data []byte) int64

	GetSubTopicIdxSet() map[string]string
}

type BaseDdCli struct {
	Host     string
	TopicSet map[string]string //已订阅的topic索引
	State    int
}
