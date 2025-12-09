package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Chandra5468/movie-streaming/database"
	"github.com/Chandra5468/movie-streaming/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

func GetMovies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30) // Use this in middleware layer
	defer cancel()

	var movies []models.Movie

	curr, err := movieCollection.Find(ctx, bson.D{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch movies"})
	}
	defer curr.Close(ctx)

	if err := curr.All(ctx, &movies); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch movies"})
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&movies)
}
