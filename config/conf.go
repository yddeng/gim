package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Address         string `toml:"Address"`
	MaxBackups      int    `toml:"MaxBackups"`
	MaxMessageCount int    `toml:"MaxMessageCount"`
	DBConfig        struct {
		SqlType  string `toml:"SqlType"`
		Host     string `toml:"Host"`
		Port     int    `toml:"Port"`
		User     string `toml:"User"`
		Password string `toml:"Password"`
		Database string `toml:"Database"`
	} `toml:"DBConfig"`
	LogConfig struct {
		Path         string `toml:"Path"`
		Debug        bool   `toml:"Debug"`
		MaxSize      int    `toml:"MaxSize"`
		EnableStdout bool   `toml:"EnableStdout"`
	} `toml:"LogConfig"`
}

var config *Config

func LoadConfig(path string) *Config {
	conf := &Config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		panic(err)
	}
	config = conf
	return conf
}

func GetConfig() *Config {
	return config
}
