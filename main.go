package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//start log

	//service initialization
	//start tcp servers
	//start quic servers
	//start http services

	//start data flow compute

	//start data dispatch services

	//start mq

	//wait for close signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
}
