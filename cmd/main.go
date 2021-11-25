package main

import (
	"github.com/yddeng/gim/config"
	_ "github.com/yddeng/gim/pkg/conv"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/utils/log"
	"os"
)

func initLog(conf *config.Config) {
	if !conf.LogConfig.Debug {
		log.CloseDebug()
	}
	if !conf.LogConfig.EnableStdout {
		log.CloseStdOut()
	}

	//log.SetOutput(conf.LogConfig.Path, conf.LogConfig.Filename, conf.LogConfig.MaxSize*1024*1024)
}

func main() {
	configPath := os.Args[1]
	conf := config.LoadConfig(configPath)

	initLog(conf)

	go func() {
		gate.StartTCPGateway("127.0.0.1:43210")
	}()

	select {}
}
