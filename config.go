package main

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

type cfgRedis struct {
	Addr string `koanf:"addr"`
}

type cfgApp struct {
	PrivateKeyPath string `koanf:"private_key_path"`
}

type Config struct {
	Redis cfgRedis
	App   cfgApp
}

var cfg Config
var k = koanf.New(".")

func initConfig() {
	// Configuration file path.
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatalf("error reading config: %v.", err)
	}

	k.Unmarshal("redis", &cfg.Redis)
	k.Unmarshal("app", &cfg.App)
}
