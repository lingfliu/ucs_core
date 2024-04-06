package mq

type InlineMq struct {
	stream_buff map[string][]byte
	streams     map[string]chan []byte
	blocks      map[string]chan string
}

func (m *InlineMq) Subscribe(topic string) chan string {
	return m.blocks[topic]
}

func (m *InlineMq) Publish(topic string, msg string) {
	m.blocks[topic] <- msg
}

func (m *InlineMq) Pull(topic string) chan []byte {
	return m.stream[topic]
}

func (m *InlineMq) Push(topic string, bs []byte) {
	m.stream_buff[topic] = append(m.stream_buff[topic], bs...)[len(bs):]
	m.streams[topic] <- m.stream_buff[topic]
}
