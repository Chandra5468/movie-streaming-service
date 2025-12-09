package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Chandra5468/movie-streaming/database"
	"github.com/Chandra5468/movie-streaming/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var validate = validator.New()

func GetMovies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30) // Use this in middleware layer
	defer cancel()

	var movies []models.Movie

	curr, err := movieCollection.Find(ctx, bson.D{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch movies"})
		return
	}
	defer curr.Close(ctx)

	if err := curr.All(ctx, &movies); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch movies"})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&movies)
}

func AddMovie(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var movie models.Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error while deserializging body" + err.Error()})
		return
	}

	// relevant validation code
	if err := validate.Struct(movie); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "error while validating body" + err.Error()})
		return
	}

	result, err := movieCollection.InsertOne(ctx, &movie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error while inserting movie" + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(result)
}
