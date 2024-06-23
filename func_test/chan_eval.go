package main

import (
	"time"

	"github.com/lingfliu/ucs_core/utils"
)

type Mq struct {
	blocks map[string]chan byte
}

var mq *Mq

func main() {
	mq = &Mq{
		blocks: make(map[string]chan byte),
	}
	go task_subscribe("mq")
	go task_publish("mq")
	tic := time.NewTicker(1 * time.Second)
	cnt := 0
	for c := range tic.C {
		cnt++
		if cnt > 10 {
			break
		}
		println(c.String())
	}
}

func task_subscribe(topic string) {
	// var b chan byte

	tic := time.NewTicker(1 * time.Second)
	for c := range tic.C {
		println(c.String())
		// b <- mq.blocks[topic]
	}
}

func task_publish(topic string) {
	idx := 1
	tick := time.NewTicker(2 * time.Millisecond)
	b := make([]byte, 1)
	for range tick.C {
		utils.Int2Byte(idx, b, 0, 1, false, true)
		mq.blocks[topic] <- b[0]
		idx++
	}
}
