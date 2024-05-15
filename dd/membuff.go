package dd

type Membuff struct {
	Buff map[string][]float64
}

func (mb Membuff) Reg(name string, buffdimen int, bufflen int) int {
	//if name is in the keys of Buff
	if _, ok := mb.Buff[name]; ok {
		mb.Buff[name] = make([]float64, bufflen)
		return 0
	} else {
		//duplicate reg, return err
		return -1
	}
}

func (mb Membuff) Unreg(name string) {
	delete(mb.Buff, name)
}

func (mb Membuff) Push(name string, data []float64) {
	//drop the oldest data
	mb.Buff[name] = append(mb.Buff[name][1:], data...)
}

/**
 * @brief merge different buff into a single vector buff
 */
func (mb Membuff) Merge(names []string) []float64 {
	d := make([]float64, len(names))
	for idx := range names {
		d[idx] = mb.Buff[names[idx]][0]
	}
	return d
}
