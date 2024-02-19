package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type IConfig interface {
	ParseFlags()
	ParseEnvs()
}

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.ServerAddress, "a", ":8080", "Server address")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "BaseURL address")
	flag.Parse()
}

func (c *Config) ParseEnvs() {
	err := env.Parse(c)
	if err != nil {
		panic(err)
	}
}

func InitConfig() *Config {
	config := &Config{}
	config.ParseFlags()
	config.ParseEnvs()
	return config
}
