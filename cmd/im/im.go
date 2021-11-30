package main

import (
	"flag"
	"fmt"
	"github.com/yddeng/gim/im"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("usage: im config")
		return
	}

	im.Service(flag.Args()[0])

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
		im.Stop()
	}
}
