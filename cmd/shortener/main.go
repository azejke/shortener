package main

import (
	"github.com/azejke/shortener/internal/handlers"
	"github.com/azejke/shortener/internal/store"
	"net/http"
)

var globalStore store.Store = make(map[string]string)

func main() {
	err := http.ListenAndServe(`:8080`, handlers.RoutesBuilder(globalStore))
	if err != nil {
		panic(err)
	}
}
