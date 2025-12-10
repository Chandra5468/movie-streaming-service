package custommiddleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func JsonRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC: %v\n%s", rec, debug.Stack())

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "internal server error", "message": "Something went wrong"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
