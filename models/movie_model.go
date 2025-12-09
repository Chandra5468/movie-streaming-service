package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required, min=2, max=100"`
}
type Ranking struct {
	RankingValue int `bson:"ranking_value" json:"ranking_value" validate:"required"`
	// RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"oneof=Excellent Good Okay Bad Terrible"`
	RankingName string `bson:"ranking_name" json:"ranking_name" validate:"required"`
}

type Movie struct {
	ID          bson.ObjectID `bson:"_id" json:"id,omitempty"`
	ImdbID      string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title       string        `bson:"title" json:"title" validate:"required,min=2,max=500"`
	PosterPath  string        `bson:"poster_path" json:"poster_path" validate:"required,url"` // moviedb website available posters can be used
	YoutubeID   string        `bson:"youtube_id" json:"youtube_id" validate:"required"`       // This id will point to trailer of each movie
	Genre       []Genre       `bson:"genres" json:"genres" validate:"required,dive"`          // keyword dive ensures nested keyword genre is also validated
	AdminReview string        `bson:"admin_review" json:"admin_review" validate:"required"`
	Ranking     Ranking       `bson:"rankings" json:"rankings" validate:"required"`
}
