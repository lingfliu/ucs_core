package conn

import (
	"github.com/lingfliu/ucs_core/coder"
)

const (
	CLI_STATE_DISCONNECTED = 0
	CLI_STATE_CONNECTING   = 1
	CLI_STATE_CONNECTED    = 2
	CLI_STATE_AUTH         = 3
)

type CliAttr struct {
	name  string
	value any
}

type Cli struct {
	Conn     Conn
	Codebook *coder.Codebook
	Coder    *coder.UCoder
	State    int

	ReqSet  map[string]coder.Msg // a req msg set that requiring response
	Auth    string
	AttrSet map[string]CliAttr

	//event
	OnReq func(coder.Msg)
}

func NewCli(connCfg *ConnCfg, codebookJson string) *Cli {
	codebook := coder.NewCodebookFromJson(codebookJson)
	var c Conn
	switch connCfg.Class {
	case CONN_CLASS_TCP:
		c = NewTcpConn(connCfg)
	case CONN_CLASS_UDP:
		c = NewUdpConn(connCfg)
	case CONN_CLASS_QUIC:
		c = NewQuicConn(connCfg)
	default:
		c = NewTcpConn(connCfg)
	}
	cli := &Cli{
		Conn:     c,
		Codebook: codebook,
		Coder:    coder.NewUCoderFromCodebook(codebook),
		State:    CLI_STATE_DISCONNECTED,
		ReqSet:   make(map[string]coder.Msg),
	}
	return cli
}

func (cli *Cli) Connect() {
	cli.State = CLI_STATE_CONNECTING
	ret := cli.Conn.Connect()
	if ret > 0 {
		cli.State = CLI_STATE_CONNECTED
	}

	// ctx, cancel := context.WithCancel(context.Background())

	// bsCh := make(chan []byte)
	// go cli.Conn.StartRecv(bsCh)
	// rxMsgChan := cli.Coder.StartDecode(bsCh)
	// go cli.HandleMsg(rxMsgChan)
}

func (cli *Cli) Disconnect() int {
	cli.Coder.StopDecode()
	return cli.Conn.Disconnect()
}

func (cli *Cli) Establish(c Conn) {
	cli.Conn = c
	cli.State = CLI_STATE_CONNECTED
}

func (cli *Cli) HandleMsg(msg chan *coder.UMsg) {
	for cli.State == CLI_STATE_CONNECTED || cli.State == CLI_STATE_AUTH {
		select {
		// case umsg := <-msg:
		//handle msg
		// msgMeta := cli.Codebook.MsgMeta[umsg.Class]
		// if msgMeta.IsExposed {
		// 	if cli.OnMsg != nil {
		// 		cli.OnMsg(umsg)
		// 	}
		// } else {
		// 	ackMsg := cli.Query(ackMsg)
		// 	bs := cli.Coder.Encode(ackMsg)
		// 	cli.ConnCli.ScheduleWrite(bs)
		// }
		}
	}
}

func (cli *Cli) Query(msg coder.Msg) *coder.UMsg {
	//flow handling
	//expose query interface
	var ackMsg *coder.UMsg = nil
	// res := cli.Codebook.GetAck(msg.GetClass())

	// if _, ok := cli.AttrSet[res]; ok {
	// 	//return attr with msg specified in rcodebook
	// 	ackMsg = cli.Coder.EncodeAttr(res, cli.AttrSet[res])
	// } else {
	// 	if cli.OnReq != nil {
	// 		ackMsg := cli.OnReq(msg, res)
	// 	}
	// }

	return ackMsg
}
