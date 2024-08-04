package types

type Node struct {
	Name    string
	Descrip string
	Attrs   map[string][]float64
	Data    map[string][]float64
	Gw      *NodeGw
}
