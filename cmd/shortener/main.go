package main

import (
	"github.com/azejke/shortener/internal/handlers"
	"github.com/azejke/shortener/internal/store"
	"net/http"
)

var Store store.Store = make(map[string]string)

func main() {
	err := http.ListenAndServe(`:8080`, http.HandlerFunc(Webhook))
	if err != nil {
		panic(err)
	}
}

func Webhook(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		handlers.GetURL(res, req, Store)
	case http.MethodPost:
		handlers.WriteURL(res, req, Store)
	default:
		res.WriteHeader(http.StatusBadRequest)
	}
}
