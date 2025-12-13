package custommiddleware

import (
	"context"
	"net/http"

	"github.com/Chandra5468/movie-streaming/utils"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// stringArray := strings.Split(authHeader, "Bearer ")
		// tokenString := stringArray[1]
		tokenString := authHeader[len("Bearer "):]

		if tokenString == "" {
			http.Error(w, "Unauthorized bearer token is required", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateToken(tokenString)

		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.UserID, claims.UserId)
		ctx = context.WithValue(ctx, utils.Role, claims.Role)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
