package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

var DefaultConfig Config

func InitConfig() *Config {
	DefaultConfig = Config{}
	flag.StringVar(&DefaultConfig.ServerAddress, "a", ":8080", "Server address")
	flag.StringVar(&DefaultConfig.BaseURL, "b", "http://localhost:8080", "BaseURL address")
	flag.Parse()
	err := env.Parse(&DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}
	return &DefaultConfig
}
