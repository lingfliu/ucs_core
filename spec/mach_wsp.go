package spec

import (
	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/conn"
)

type MachWspCoder struct {
}

type MachWspConnCli struct {
	conn.Cli
}

func (m *MachWspConnCli) Test() {
}

func (m *MachWspConnCli) HandleMsg(msg *coder.UMsg) {
}
