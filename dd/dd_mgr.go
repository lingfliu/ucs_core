package dd

type DdMgr struct {
	dds map[string]*Dd
}

func NewDdMgr() *DdMgr {
	return &DdMgr{
		dds: make(map[string]*Dd),
	}
}

func (ddm *DdMgr) RegDd(name string, dd *Dd) {
	ddm.dds[name] = dd
}

func (ddm *DdMgr) UnregDd(name string) {
	delete(ddm.dds, name)
}
