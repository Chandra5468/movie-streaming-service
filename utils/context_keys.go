package utils

import (
	"errors"
	"net/http"
)

type ContextKey string

const (
	UserID ContextKey = "userId"
	Role   ContextKey = "role"
	Email  ContextKey = "email"
)

// Even better use a struct with combination of above consts
// And keep this file in types than utils

func GetDataFromContext(r *http.Request) (ContextKey, error) {
	userId := r.Context().Value(UserID)
	if userId == nil {
		return "", errors.New("userid does not exists in context")
	}

	id, ok := userId.(ContextKey)

	if !ok {
		return "", errors.New("unable to retrive userid")
	}

	return id, nil
}
