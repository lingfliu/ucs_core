package types

type Pos struct {
	// GNSS position
	longitude float64
	latitude  float64
	altitude  float64

	// local position
	x        float64
	y        float64
	h        float64
	addr_maj float64
	addr_min float64
}
