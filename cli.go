package main

import (
	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/conn"
)

const (
	CLI_STATE_DISCONNECTED = 0
	CLI_STATE_CONNECTING   = 1
	CLI_STATE_CONNECTED    = 2
	CLI_STATE_AUTH         = 3
)

type Cli struct {
	ConnCli  *conn.TcpConnCli
	Codebook *coder.Codebook
	Coder    *coder.Coder
	State    int

	FlowState int
	ReqSet    map[string]coder.Msg
}

func NewCli(connCfg *conn.ConnCfg, codebook *coder.Codebook) *Cli {
	cli := &Cli{
		ConnCli:  conn.NewTcpConnCli(connCfg.RemoteAddr, connCfg.Port, true, 1000),
		Codebook: codebook,
		Coder:    coder.NewCodebookFromJson(codebook),
		State:    0,
		ReqSet:   make(map[string]coder.Msg),
	}
	return cli
}

func (cli *Cli) Connect() int {
	return cli.ConnCli.Connect()
}

func (cli *Cli) Disconnect() int {
	return cli.ConnCli.Disconnect()
}

func (cli *Cli) HandleMsg(msg coder.UMsg) {
}

func (cli *Cli) Req(msg coder.Msg) *coder.UMsg {

	return &coder.UMsg{}
}
