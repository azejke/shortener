package main

import (
	"github.com/azejke/shortener/internal/app"
	"github.com/azejke/shortener/internal/store"
	"net/http"
)

var Store store.Store = make(map[string]string)

func main() {
	err := http.ListenAndServe(`:8080`, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.Run(writer, request, Store)
	}))
	if err != nil {
		panic(err)
	}
}
