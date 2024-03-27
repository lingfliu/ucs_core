package types

type NodeGw struct {
	// connMeta *ConnMeta
	id  string
	mac string

	connType string

	nodes map[string]*Node
}
