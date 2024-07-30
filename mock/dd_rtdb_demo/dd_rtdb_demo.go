package main

import (
	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/dd"
)

func main() {
	rtdb := &rtdb.AgilorRtdb{}

	go rtdb.Connect()

	memDd := &dd.MemDd{}

	memDd.Listen()

}
