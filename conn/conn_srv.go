package conn

type ConnSrv struct {
	Port int
	KeepAlive bool
	ConnCliSet map[string]*ConnCli
	State int
}

func (srv *ConNSrv) _task_cleanup() {
	tic := time.NewTicker(time.Second * 1)
	for srv.State == STATE_START {
		select {
			case <-tic.C:
				for _, v := range srv.TcpConnCliSet {
				}
	}
}

