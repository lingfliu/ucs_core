package rtdb

// #cgo LDFLAGS: -lagilor
// #include "./include/agilor_defs.h"
// #include "./include/agilor.h"
import "C"

const (
	AGILOR_STATE_DISCONNECTED = 0
	AGILOR_STATE_CONNECTED    = 1
)

type AgilorCli struct {
	SrvAddr  string
	SrvPort  int
	Username string
	Password string
	State    int
}

func NewAgilorCli(srvAddr string, srvPort int, username string, password string) *AgilorCli {
	return &AgilorCli{
		SrvAddr:  srvAddr,
		SrvPort:  srvPort,
		Username: username,
		Password: password,
		State:    AGILOR_STATE_DISCONNECTED,
	}
}

func (cli *AgilorCli) Connect() {
	//TODO: connect using C API
}

func (cli *AgilorCli) Disconnect() {
}

func (cli *AgilorCli) Update() {
}
