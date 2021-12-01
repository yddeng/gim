package im

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Address          string `toml:"Address"`
	MaxBackups       int    `toml:"MaxBackups"`
	MaxMessageCount  int    `toml:"MaxMessageCount"`
	UserCacheCount   int    `toml:"UserCacheCount"`
	GroupCacheCount  int    `toml:"GroupCacheCount"`
	MaxTaskCount     int    `toml:"MaxTaskCount"`
	HeartbeatTimeout int    `toml:"HeartbeatTimeout"`
	DBConfig         struct {
		SqlType  string `toml:"SqlType"`
		Host     string `toml:"Host"`
		Port     int    `toml:"Port"`
		User     string `toml:"User"`
		Password string `toml:"Password"`
		Database string `toml:"Database"`
	} `toml:"DBConfig"`
	LogConfig struct {
		Path         string `toml:"Path"`
		Filename     string `toml:"Filename"`
		Debug        bool   `toml:"Debug"`
		MaxSize      int    `toml:"MaxSize"`
		EnableStdout bool   `toml:"EnableStdout"`
	} `toml:"LogConfig"`
}

func loadCfg(path string) *Config {
	conf := &Config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		panic(err)
	}
	return conf
}

var config *Config
