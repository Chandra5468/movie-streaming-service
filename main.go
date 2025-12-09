package main

import (
	"log"
	"net/http"

	"github.com/Chandra5468/movie-streaming/controllers"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/v1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello"))
	})
	router.Get("/movies", controllers.GetMovies)
	router.Post("/movie", controllers.AddMovie)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server not started error %v", err)
	}
}
