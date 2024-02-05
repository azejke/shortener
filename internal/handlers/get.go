package handlers

import (
	"github.com/azejke/shortener/internal/store"
	"net/http"
	"strings"
)

func GetURL(res http.ResponseWriter, req *http.Request, store store.Store) {
	contentTypeValue := req.Header.Get("Content-Type")
	id := strings.Trim(req.URL.Path, "/")
	url, ok := store[id]
	if contentTypeValue != "text/plain" || !ok || len(id) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
	return
}
