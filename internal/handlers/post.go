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
	body, err := io.ReadAll(req.Body)
	contentTypeValue := req.Header.Get("Content-Type")
	log.Printf("Content-Type value: %s", contentTypeValue)
	if err != nil || len(body) == 0 || contentTypeValue != "text/plain" {
		res.WriteHeader(http.StatusBadRequest)
	}
	generatedKey := generateRandomString(10)
	store[generatedKey] = string(body)
	host := req.Host
	log.Printf("Host: %s", host)
	res.Header().Set(`Content-Type`, `text/plain`)
	res.WriteHeader(http.StatusCreated)
	result := fmt.Sprintf("http://%s/%s", host, generatedKey)
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
