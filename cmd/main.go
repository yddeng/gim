package main

import (
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/logger"
)

func main() {
	logger.InitLogger("log", "gim.log")

	go func() {
		gate.StartTCPGateway("127.0.0.1:43210")
	}()

}
