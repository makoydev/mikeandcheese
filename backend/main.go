package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/makoydev/mikeandcheese/backend/internal/database"
	"github.com/makoydev/mikeandcheese/backend/internal/handlers"
)

func main() {
	// Initialize database
	db, err := database.Connect("blog.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.CreateTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	if err := database.SeedData(db); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	// Set up router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	postHandler := handlers.NewPostHandler(db)

	r.Route("/api", func(r chi.Router) {
		r.Get("/posts", postHandler.GetPosts)
		r.Get("/posts/featured", postHandler.GetFeaturedPosts)
		r.Get("/posts/{slug}", postHandler.GetPost)
	})

	// Server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-done
	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
