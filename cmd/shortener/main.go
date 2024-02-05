package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urls = make(map[string]string)

func main() {
	err := http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
	if err != nil {
		panic(err)
	}
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

func webhook(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		body, _ := io.ReadAll(req.Body)
		if len(string(body)) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		generatedKey := generateRandomString(10)
		urls[generatedKey] = string(body)
		host := req.Host
		res.Header().Set(`Content-Type`, `text/plain`)
		res.WriteHeader(http.StatusCreated)
		result := fmt.Sprintf("http://%s/%s", host, generatedKey)
		_, _ = res.Write([]byte(result))
	} else if req.Method == http.MethodGet {
		id := strings.Trim(req.URL.Path, "/")
		url, ok := urls[id]
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
