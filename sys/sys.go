package sys

import (
	"github.com/lingfliu/ucs_core/conn/srv"
)

type Sys struct {
	servers map[string]srv.Server
}

func NewDefaultSys() *Sys {
	return &Sys{}
}

/**
 * @Description: register new server
 * @param tpl
 */
func (s *Sys) RegServer(name string, port int, connMode int) {
	s.servers[name] = srv.NewServer(name, port, connMode)
	s.servers[name].Start()
}

func (s *Sys) Start() {
	//
}

func (s *Sys) Stop() {
}

func (s *Sys) Shutdown() {
}

func (s *Sys) Config() {
}
