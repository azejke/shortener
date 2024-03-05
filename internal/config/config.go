package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type IConfig interface {
	parseFlags()
	parseEnvs()
}

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	LogLevel      string `env:"LOG_LEVEL"`
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.ServerAddress, "a", ":8080", "Server address")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "BaseURL address")
	flag.StringVar(&c.LogLevel, "l", "info", "Log level")
	flag.Parse()
}

func (c *Config) parseEnvs() {
	err := env.Parse(c)
	if err != nil {
		panic(err)
	}
}

func InitConfig() *Config {
	config := &Config{}
	config.parseFlags()
	config.parseEnvs()
	return config
}
