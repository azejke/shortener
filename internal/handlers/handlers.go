package handlers

import (
	"fmt"
	"github.com/azejke/shortener/config"
	"github.com/azejke/shortener/internal/store"
	"github.com/azejke/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
)

func RoutesBuilder() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		SearchURL(writer, request)
	})
	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		WriteURL(writer, request)
	})
	return r
}

func SearchURL(res http.ResponseWriter, req *http.Request) {
	s := *&store.Store
	id := chi.URLParam(req, "id")
	log.Printf("Received id: %s", id)
	urlValue, ok := s[id]
	if !ok || len(id) == 0 {
		log.Println("URL is empty or doesn't exist")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("URL value: %s", urlValue)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set("Location", urlValue)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func WriteURL(res http.ResponseWriter, req *http.Request) {
	contentTypeValue := req.Header.Get("Content-Type")
	log.Printf("Content-Type value: %s", contentTypeValue)
	if contentTypeValue != "text/plain; charset=utf-8" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(req.Body)
	log.Printf("Body: %s", string(body))
	if err != nil || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	generatedKey := utils.GenerateRandomString(10)
	s := *&store.Store
	s[generatedKey] = string(body)
	c := *&config.DefaultConfig
	log.Printf("baseURL: %s", c.BaseURL)
	res.Header().Set(`Content-Type`, `text/plain; charset=utf-8`)
	res.WriteHeader(http.StatusCreated)
	result := fmt.Sprintf("%s%s/%s", c.BaseURL, c.ServerAddress, generatedKey)
	log.Printf("Result value: %s", result)
	_, _ = res.Write([]byte(result))
}
