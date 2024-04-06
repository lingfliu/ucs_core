package srv

import (
	"github.com/lingfliu/ucs_core/types"
)

type Server interface {
	Start() chan types.Msg
	Stop() error
	Shutdown() error
}
