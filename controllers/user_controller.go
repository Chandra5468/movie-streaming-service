package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Chandra5468/movie-streaming/database"
	"github.com/Chandra5468/movie-streaming/models"
	"github.com/Chandra5468/movie-streaming/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection("users")

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "validation failed"})
		return
	}

	hashedPwd, err := HashPassword(user.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "password not stored"})
		return
	}

	count, err := userCollection.CountDocuments(r.Context(), bson.D{
		bson.E{
			Key:   "email",
			Value: user.Email,
		},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to check existing user"})
		return
	}

	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "User already exists"})
		return
	}

	user.UserID = primitive.NewObjectID().Hex()
	user.Password = hashedPwd
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = userCollection.InsertOne(r.Context(), user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "successful"})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}

	var foundUser models.User
	err := userCollection.FindOne(r.Context(), bson.D{
		bson.E{
			Key:   "email",
			Value: userLogin.Email,
		},
	}).Decode(&foundUser)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid password"})
		return
	}

	token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
		return
	}

	err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to update tokens"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.UserResponse{
		UserId:          foundUser.UserID,
		FirstName:       foundUser.FirstName,
		LastName:        foundUser.LastName,
		Email:           foundUser.Email,
		Role:            foundUser.Role,
		Token:           token,
		RefreshToken:    refreshToken,
		FavouriteGenres: foundUser.FavouriteGenres,
	})

}
