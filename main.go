package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chandra5468/movie-streaming/controllers"
	"github.com/Chandra5468/movie-streaming/database"
	custommiddleware "github.com/Chandra5468/movie-streaming/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initializing MongoDB Client
	database.GetClient()

	// Create the router and apply middleware
	router := chi.NewRouter()
	router.Use(middleware.Logger)    // Log all HTTP requests
	router.Use(middleware.Recoverer) // Recover from panics
	// global custom middleware
	router.Use(custommiddleware.CORS)

	router.Route("/api", func(r chi.Router) {
		r.Post("/register", controllers.RegisterUser)
		r.Post("/login", controllers.LoginUser)

		// Protected routes
		r.Group(func(protected chi.Router) {
			protected.Use(custommiddleware.Auth)
			protected.Post("/movie", controllers.AddMovie)
			protected.Get("/movies", controllers.GetMovies)
			protected.Patch("/updatereview/:imdb_id", controllers.AdminReviewUpdate)
			protected.Get("/recommended/movies", controllers.GetRecommendedMovies)
		})
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// handling graceful shutdowns
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// wait for interrupt signal to gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down server...")

	// Create a context with a timeout to ensure the server shuts down properly
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}
	// Disconnect MongoDB client gracefully
	database.Disconnect()
	log.Println("Server gracefully stopped")
}
