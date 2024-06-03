package dd

type Dd struct {
	// 用于存储数据的通道
	SrcStream   chan byte
	DistStreams map[string]*DdSubStream
}

func NewDd(source chan byte) *Dd {
	return &Dd{
		SrcStream:   source,
		DistStreams: make(map[string]*DdSubStream),
	}
}

func (dd *Dd) RegStream(name string) {
	dd.DistStreams[name] = make(*DdSubStream)
}

func (dd *Dd) UnregStream(name string) {
	//remove the channel from the map
	delete(dd.SubStreams, name)
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
