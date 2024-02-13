package config

import (
	"flag"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

var DefaultConfig = &Config{}

func (c *Config) InitConfig() *Config {
	flag.StringVar(&c.ServerAddress, "a", ":8080", "Server address")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080/", "BaseURL address")
	flag.Parse()
	return c
}
