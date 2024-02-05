package main

import (
	"github.com/azejke/shortener/internal/app"
	"github.com/azejke/shortener/internal/store"
	"net/http"
)

var Store store.Store = make(map[string]string)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		app.Run(writer, request, Store)
	})
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
