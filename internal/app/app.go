package app

import (
	"github.com/azejke/shortener/internal/handlers"
	"github.com/azejke/shortener/internal/store"
	"net/http"
)

func Run(res http.ResponseWriter, req *http.Request, store store.Store) {
	webhook(res, req, store)
}

func webhook(res http.ResponseWriter, req *http.Request, store store.Store) {
	switch req.Method {
	case http.MethodGet:
		handlers.SearchURL(res, req, store)
	case http.MethodPost:
		handlers.WriteURL(res, req, store)
	default:
		res.WriteHeader(http.StatusBadRequest)
	}
}
