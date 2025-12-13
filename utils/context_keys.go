package utils

type ContextKey string

const (
	UserID ContextKey = "userId"
	Role   ContextKey = "role"
	Email  ContextKey = "email"
)

// Even better use a struct with combination of above consts
// And keep this file in types than utils
