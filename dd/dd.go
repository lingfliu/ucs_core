package dd

type Dd struct {
	// 用于存储数据的通道
	Source  chan byte
	Streams map[string]chan byte
}

func NewDd(source chan byte) *Dd {
	return &Dd{
		Source:  source,
		Streams: make(map[string]chan byte),
	}
}

func (dd *Dd) RegStream(name string) {
	dd.Streams[name] = make(chan byte)
}

func (dd *Dd) UnregStream(name string) {
	//remove the channel from the map
	delete(dd.Streams, name)
}

func (dd *Dd) Run() {
	for range dd.Source {
		select {
		case data := <-dd.Source:
			for _, stream := range dd.Streams {
				stream <- data
			}
		}

	}
}

func (dd *Dd) Close() {
	//TODO: remove all streams and close the source
	for _, stream := range dd.Streams {
		close(stream)
	}
}
