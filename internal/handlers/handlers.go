package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/azejke/shortener/internal/config"
	"github.com/azejke/shortener/internal/models"
	"github.com/azejke/shortener/internal/store"
	"github.com/azejke/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type IURLHandler interface {
	SearchURL(res http.ResponseWriter, req *http.Request)
	WriteURL(res http.ResponseWriter, req *http.Request, cfg *config.Config)
	Shorten(res http.ResponseWriter, req *http.Request, cfg *config.Config)
}

type URLHandler struct {
	storage *store.Store
	cfg     *config.Config
}

func (u *URLHandler) SearchURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	urlValue, ok := u.storage.Get(id)
	if !ok || len(id) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Println(urlValue)
	res.Header().Set("Location", urlValue)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (u *URLHandler) WriteURL(res http.ResponseWriter, req *http.Request) {
	contentTypeValue := req.Header.Get("Content-Type")
	if contentTypeValue == "text/plain; charset=utf-8" || contentTypeValue == "application/x-gzip" {
		body, err := io.ReadAll(req.Body)
		if err != nil || len(body) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		generatedKey := utils.GenerateRandomString(10)
		u.storage.Insert(generatedKey, string(body))
		res.Header().Set(`Content-Type`, `text/plain; charset=utf-8`)
		res.WriteHeader(http.StatusCreated)
		result := fmt.Sprintf("%s/%s", u.cfg.BaseURL, generatedKey)
		_, _ = res.Write([]byte(result))
		return
	}
	res.WriteHeader(http.StatusBadRequest)
}

func (u *URLHandler) Shorten(res http.ResponseWriter, req *http.Request) {
	contentTypeValue := req.Header.Get("Content-Type")
	if contentTypeValue != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var url models.ShortenRequest
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &url); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if url.URL == "" {
		http.Error(res, "Invalid request params", http.StatusInternalServerError)
		return
	}
	var result models.ShortenResponse
	generatedKey := utils.GenerateRandomString(10)
	u.storage.Insert(generatedKey, url.URL)
	result.Result = fmt.Sprintf("%s/%s", u.cfg.BaseURL, generatedKey)
	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(res, "Encoding error", http.StatusInternalServerError)
		return
	}
	res.Header().Set(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}
