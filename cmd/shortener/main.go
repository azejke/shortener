package main

import (
	"github.com/azejke/shortener/internal/app"
	"github.com/azejke/shortener/internal/store"
	"net/http"
)

var globalStore store.Store = make(map[string]string)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		app.Run(writer, request, globalStore)
	})
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
