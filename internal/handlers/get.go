package handlers

import (
	"github.com/azejke/shortener/internal/store"
	"log"
	"net/http"
	"strings"
)

func SearchURL(res http.ResponseWriter, req *http.Request, store store.Store) {
	contentTypeValue := req.Header.Get("Content-Type")
	if contentTypeValue != "text/plain; charset=utf-8" {
		log.Printf("Incorrect Content-Type: %s", contentTypeValue)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	id := strings.Trim(req.URL.Path, "/")
	log.Printf("Received id: %s", id)
	urlValue, ok := store[id]
	if !ok || len(id) == 0 {
		log.Println("URL is empty or doesn't exist")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("URL value: %s", urlValue)
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", urlValue)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
