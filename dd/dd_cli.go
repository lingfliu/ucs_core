package dd

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
	Host           string
	SubTopicIdxSet map[string]string //已订阅的topic索引
}
