package prop

type Range struct {
	//if both are 0, no range check
	Min float64 //by default, 0
	Max float64 //by default, 0
}

func (r *Range) IsOverflow(value float64) bool {
	return value >= r.Min && value <= r.Max
}
