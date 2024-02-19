package main

import (
	"github.com/azejke/shortener/internal/config"
	"github.com/azejke/shortener/internal/handlers"
	"github.com/azejke/shortener/internal/store"
	"net/http"
)

func main() {
	cfg := config.InitConfig()
	storage := store.InitStore()
	err := http.ListenAndServe(cfg.ServerAddress, handlers.RoutesBuilder(cfg, storage))
	if err != nil {
		panic(err)
	}
}
