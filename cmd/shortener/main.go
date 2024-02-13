package main

import (
	"github.com/azejke/shortener/config"
	"github.com/azejke/shortener/internal/handlers"
	"net/http"
)

func main() {
	c := config.DefaultConfig.InitConfig()
	err := http.ListenAndServe(c.ServerAddress, handlers.RoutesBuilder())
	if err != nil {
		panic(err)
	}
}
