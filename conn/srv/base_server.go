package srv

import (
	"github.com/lingfliu/ucs_core/conn/conn"
)

type BaseServer struct {
	Name       string
	Port       int
	ConnMode   int
	ListenConn conn.Conn

	CliConn map[string]conn.Conn
}

func (s *BaseServer) Start() error {
	go s.task_cleanup()
	return nil
}

func (s *BaseServer) Stop() error {
	return nil
}

func (s *BaseServer) Shutdown() error {
	return nil
}

func (s *BaseServer) task_cleanup() {
	// for n := range s.CliConn {
	// c := s.CliConn[n]
	// if !c.IsActive() {
	// s.CloseCli(c)
	// }
	// }
}
