package dd

type DataDispatch struct {
	// 用于存储数据的通道
	source       chan byte
	dispatchChan chan map[string]chan byte
	membuff      Membuff
}

func NewDataDispatch(source chan byte) *DataDispatch {
	return &DataDispatch{
		source:       source,
		dispatchChan: make(chan map[string]chan byte),
	}
}

func (dd *DataDispatch) Reg(names []string, ch chan byte) {
	// dd.dispatchChan <- map[string]chan byte
}

func (dd *DataDispatch) Unreg(name string) {
	//remove the channel from the map
	dd.dispatchChan <- map[string]chan byte{name: nil}
}

func (dd *DataDispatch) Run() {
	for {
		select {
		case data := <-dd.source:
			for _, ch := range <-dd.dispatchChan {
				ch <- data
			}
		}
	}
}

func (dd *DataDispatch) Close() {
	close(dd.source)
	close(dd.dispatchChan)
}
