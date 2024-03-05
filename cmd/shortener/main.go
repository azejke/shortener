package main

import (
	"github.com/azejke/shortener/internal/config"
	"github.com/azejke/shortener/internal/handlers"
	"github.com/azejke/shortener/internal/logger"
	"github.com/azejke/shortener/internal/store"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg := config.InitConfig()
	storage := store.InitStore()

	if err := logger.InitLogger(cfg.LogLevel); err != nil {
		panic(err)
	}

	logger.Log.Info("Running server", zap.String("address", cfg.ServerAddress))

	err := http.ListenAndServe(cfg.ServerAddress, handlers.RoutesBuilder(cfg, storage))
	if err != nil {
		panic(err)
	}
}
