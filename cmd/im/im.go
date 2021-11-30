package main

import (
	"flag"
	"fmt"
	"github.com/yddeng/gim/im"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("usage: im config")
		return
	}

	im.Service(flag.Args()[0])
	select {}
}
