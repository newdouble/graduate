package config

import (
	"gopkg.in/gcfg.v1"
)

var defaultFilePath = "./conf/server.conf"

var gConf Config

type Config struct {
	ServerConf ServerConfig
	LocationConf LocationConfig
	SwitchConf SwitchConfig
}

type SwitchConfig struct {
	ManageActMethodAllow []string
}

type ServerConfig struct {
	Port string
	ReadTimeout string
	WriteTimeout string
}

type LocationConfig struct {
	ManageActUpdatePath string
	GrantRightPath  string
	FetchActivityDetailPath string
}


func Get() *Config {
	return &gConf
}

func Parse(fp string) {
	if fp == "" {
		fp = defaultFilePath
	}

	if err := gcfg.ReadFileInto(&gConf, fp); err != nil {
		panic(err)
	}
}