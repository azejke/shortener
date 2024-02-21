package handlers

import (
	"github.com/azejke/shortener/internal/config"
	"github.com/azejke/shortener/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func RoutesBuilder(cfg *config.Config, s *store.Store) chi.Router {
	handlers := URLHandler{storage: s}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{id}", handlers.SearchURL)
	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.WriteURL(writer, request, cfg)
	})
	return r
}
