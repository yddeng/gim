package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Address         string `toml:"Address"`
	MaxMessageCount int    `toml:"MaxMessageCount"`
	DBConfig        struct {
		SqlType  string `toml:"SqlType"`
		Host     string `toml:"Host"`
		Port     int    `toml:"Port"`
		User     string `toml:"User"`
		Password string `toml:"Password"`
		DataBase string `toml:"DataBase"`
	} `toml:"DBConfig"`
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
