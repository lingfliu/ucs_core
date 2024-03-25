package stats

func Minmax(dao Dao, range Range) (int, int) {
	min := range[0]
	max := range[0]
	for _, v := range {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}
