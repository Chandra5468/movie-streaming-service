package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
}
