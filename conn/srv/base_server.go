package srv

import (
	"github.com/lingfliu/ucs_core/conn/conn"
	"github.com/lingfliu/ucs_core/types"
)

func NewServer(name string, port int, connMode int) Server {
	return &BaseServer{
		Name:     name,
		Port:     port,
		ConnMode: connMode,
	}
}

type BaseServer struct {
	Name     string
	Port     int
	ConnMode int
	Conn     conn.Conn
}

func (s *BaseServer) Start() chan types.Msg {
	return make(chan types.Msg)
}

func (s *BaseServer) Stop() error {
	return nil
}

func (s *BaseServer) Shutdown() error {
	return nil
}

func (s *BaseServer) OnMsg(msg types.Msg) {
}
