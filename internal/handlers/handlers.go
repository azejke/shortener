package handlers

import (
	"fmt"
	"github.com/azejke/shortener/internal/config"
	"github.com/azejke/shortener/internal/store"
	"github.com/azejke/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

type IURLHandler interface {
	SearchURL(res http.ResponseWriter, req *http.Request)
	WriteURL(res http.ResponseWriter, req *http.Request, cfg *config.Config)
}

type URLHandler struct {
	storage *store.Store
}

func (u *URLHandler) SearchURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	urlValue, ok := u.storage.Get(id)
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

func (u *URLHandler) WriteURL(res http.ResponseWriter, req *http.Request, cfg *config.Config) {
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
	u.storage.Insert(generatedKey, string(body))
	log.Printf("BaseURL: %s", cfg.BaseURL)
	res.Header().Set(`Content-Type`, `text/plain; charset=utf-8`)
	res.WriteHeader(http.StatusCreated)
	result := fmt.Sprintf("%s/%s", cfg.BaseURL, generatedKey)
	log.Printf("Result value: %s", result)
	_, _ = res.Write([]byte(result))
}
