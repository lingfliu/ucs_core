package main

import (
	"github.com/lingfliu/ucs_core/data/rtdb"
)

func main() {
	rtdb := &rtdb.AgilorRtdb{}

	go rtdb.Connect()

}
