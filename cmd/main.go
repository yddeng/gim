package main

import (
	_ "github.com/yddeng/gim/pkg/chat"
	"github.com/yddeng/gim/pkg/gate"
)

func main() {
	//logger.InitLogger("log", "gim.log")

	go func() {
		gate.StartTCPGateway("127.0.0.1:43210")
	}()

	select {}
}