package handlers

import (
	"fmt"
	"github.com/azejke/shortener/internal/store"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func WriteURL(res http.ResponseWriter, req *http.Request, store store.Store) {
	contentTypeValue := req.Header.Get("Content-Type")
	if contentTypeValue != "text/plain; charset=utf-8" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(req.Body)
	log.Printf("Content-Type value: %s", contentTypeValue)
	if err != nil || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	generatedKey := generateRandomString(10)
	store[generatedKey] = string(body)
	host := req.Host
	log.Printf("Host: %s", host)
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	res.Header().Set(`Content-Type`, `text/plain`)
	res.WriteHeader(http.StatusCreated)
	result := fmt.Sprintf("%s://%s/%s", scheme, host, generatedKey)
	log.Printf("Result value: %s", result)
	_, _ = res.Write([]byte(result))
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}
