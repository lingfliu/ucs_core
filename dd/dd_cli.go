package dd

type DdCli struct {
}

func (cli *DdCli) _connect(res chan int) {
	//connect coroutine
	res <- 0
}

func (cli *DdCli) Connect(on_connect func()) {
	cli._connect(make(chan int))
}

func (cli *DdCli) Disconnect() {

}

func (cli *DdCli) Subscribe(topic string, on_msg func(json string)) {
	//subscribe with json call back
}

func (cli *DdCli) Unsubscribe(topic string) {

}

func (cli *DdCli) Publish(topic string, json string) {

}
