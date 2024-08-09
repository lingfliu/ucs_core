package main

import (
	"context"
	"strconv"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
)

func task_1(ctx context.Context) {
	cnt := 0
	defer ulog.Log().I("task_1", "task_2 done")
	for {
		select {
		// case <-ctx.Done():
		// return
		case stop := <-ctx.Value("sig_stop").(chan bool):
			if stop {
				return
			}
		default:
			// do something
			time.Sleep(1 * time.Second)
			ulog.Log().I("task_1", "doing something"+" "+strconv.Itoa(cnt))
			cnt += 1
		}
	}

}

func task_2(ctx context.Context) {
	cnt := 0
	defer ulog.Log().I("task_2", "task_2 done")
	for {
		select {
		// case <-ctx.Done():
		// return
		case stop := <-ctx.Value("sig_stop").(chan bool):
			if stop {
				return
			}
		default:
			// do something
			time.Sleep(1 * time.Second)
			ulog.Log().I("task_2", "doing something"+" "+strconv.Itoa(cnt))
			cnt += 1
		}
	}

}
func main() {
	ulog.Config(ulog.LOG_LEVEL_INFO, "", false)
	ctx, cancel := context.WithCancel(context.Background())
	// ctx := context.WithValue(context.Background(), "sig_stop", make(chan bool))

	go task_1(ctx)
	go task_2(ctx)

	tic := time.Tick(5 * time.Second)
	select {
	case <-tic:
		// ctx.Value("sig_stop").(chan bool) <- true
		cancel()
		break
	}

	for {
		time.Sleep(1 * time.Second)
	}

}
