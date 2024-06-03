package dd

type DdSubStream struct {
	KaStream chan byte // keep alive stream for sending ping-pong message
	Stream   chan byte // data stream
}

func (dss *DdSubStream) Close() {
	close(dss.KaStream)
	close(dss.Stream)
}

func (dss *DdSubStream) Send(data byte) {
	dss.Stream <- data
}

func (dss *DdSubStream) _task_ka() {
	for range dss.KaStream {
		dss.KaStream <- 0
	}
}
