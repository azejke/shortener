package handlers

import (
	"compress/gzip"
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
	WriteURL(res http.ResponseWriter, req *http.Request)
	Shorten(res http.ResponseWriter, req *http.Request)
	generateAndInsertKey(url string) string
	checkContentType(contentType string, expectedContentType ...string) bool
}

type URLHandler struct {
	storage *store.Store
	cfg     *config.Config
}

func (u *URLHandler) checkContentType(actualContentType string, expectedContentType ...string) bool {
	isExist := false
	for _, v := range expectedContentType {
		if v == actualContentType {
			isExist = true
		}
	}
	return isExist
}

func (u *URLHandler) generateAndInsertKey(url string) string {
	generatedKey := utils.GenerateRandomString(10)
	u.storage.Insert(generatedKey, url)
	return generatedKey
}

func (u *URLHandler) SearchURL(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	urlValue, ok := u.storage.Get(id)
	if !ok || len(id) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set("Location", urlValue)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (u *URLHandler) WriteURL(res http.ResponseWriter, req *http.Request) {
	contentTypeValue := req.Header.Get("Content-Type")
	if !u.checkContentType(contentTypeValue, "text/plain; charset=utf-8", "application/x-gzip") {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var body []byte
	if contentTypeValue == "text/plain; charset=utf-8" {
		body, _ = io.ReadAll(req.Body)
		if len(body) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		body, err = io.ReadAll(gz)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	key := u.generateAndInsertKey(string(body))
	res.Header().Set(`Content-Type`, `text/plain; charset=utf-8`)
	res.WriteHeader(http.StatusCreated)
	result := fmt.Sprintf("%s/%s", u.cfg.BaseURL, key)
	_, _ = res.Write([]byte(result))
}

func (u *URLHandler) Shorten(res http.ResponseWriter, req *http.Request) {
	contentTypeValue := req.Header.Get("Content-Type")
	if !u.checkContentType(contentTypeValue, "application/json") {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var sreq *models.ShortenRequest
	if err := json.NewDecoder(req.Body).Decode(&sreq); err != nil {
		http.Error(res, "Cant read body", http.StatusBadRequest)
		return
	}
	originalURL := sreq.URL
	if originalURL == "" {
		http.Error(res, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	var sres models.ShortenResponse
	key := u.generateAndInsertKey(originalURL)
	sres.Result = fmt.Sprintf("%s/%s", u.cfg.BaseURL, key)
	res.Header().Set(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(res).Encode(sres); err != nil {
		http.Error(res, "error encoding response", http.StatusInternalServerError)
		return
	}
}
