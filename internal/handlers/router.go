package handlers

import (
	"github.com/azejke/shortener/internal/config"
	"github.com/azejke/shortener/internal/logger"
	"github.com/azejke/shortener/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var defaultCompressibleContentTypes = []string{
	"text/html",
	"text/css",
	"text/plain",
	"text/plain; charset=utf-8",
	"text/xml",
	"text/javascript",
	"application/javascript",
	"application/json",
}

func RoutesBuilder(cfg *config.Config, s *store.Store) chi.Router {
	handlers := URLHandler{storage: s, cfg: cfg}
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	//r.Use(middlewares.GzipHandle)
	r.Use(middleware.Compress(5, defaultCompressibleContentTypes...))
	r.Post("/", handlers.WriteURL)
	r.Get("/{id}", handlers.SearchURL)
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handlers.Shorten)
	})
	return r
}
