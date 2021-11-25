package main

import (
	"github.com/yddeng/gim/config"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/gim/pkg/conv"
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

	if err := db.Open(conf.DBConfig.SqlType,
		conf.DBConfig.Host,
		conf.DBConfig.Port,
		conf.DBConfig.Database,
		conf.DBConfig.User,
		conf.DBConfig.Password); err != nil {
		panic(err)
	}

	go func() {
		gate.StartTCPGateway("127.0.0.1:43210")
	}()

	conv.InitMessageTable()

	select {}
}
