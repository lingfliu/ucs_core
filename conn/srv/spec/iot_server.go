package spec

import (
	"github.com/lingfliu/ucs_core/conn/coder"
	"github.com/lingfliu/ucs_core/conn/srv"
)

type IotServer struct {
	srv.BaseServer
}

func (s *IotServer) NewIotServer(connMode int, codebook coder.Codebook, port int) {
	//
}
