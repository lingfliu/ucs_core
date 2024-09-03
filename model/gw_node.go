package model

import "github.com/lingfliu/ucs_core/conn"

/**
 * 网关节点，一般情况下仅用于进行数据转发
 */
type GwNode struct {
	Id         string
	Mac        string
	Addr       string //ip or url
	Conn       conn.Conn
	DNodeSet   map[string]*DNode
	CtlNodeSet map[string]*CtlNode
}
