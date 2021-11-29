package main

import (
	"github.com/yddeng/gim/im"
	"os"
)

func main() {
	im.Service(os.Args[1])
	select {}
}
