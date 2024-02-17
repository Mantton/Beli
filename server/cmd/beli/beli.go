package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mantton/beli/internal/cache"
	"github.com/mantton/beli/internal/env"
	v1 "github.com/mantton/beli/internal/handlers/v1"
	"github.com/rs/cors"
)

func main() {
	// Cache
	cache := cache.New()

	err := cache.Connect()

	if err != nil {
		log.Fatal(err)
	}

	// V1 Route Handler
	v1 := v1.New(cache)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(15 * time.Second))

	// Set up CORS middleware to allow all origins
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// Use the CORS middleware
	r.Use(corsMiddleware.Handler)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("B E L I\nE L I  \nL I     \nI\n"))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Post("/draw", v1.HandleDrawTile)
		r.Get("/info", v1.HandleGetTile)
		r.Get("/board", v1.HandleGetBoard)
	})

	slog.Info("Starting Server.")
	err = http.ListenAndServe(fmt.Sprintf(":%s", env.Port()), r)

	if err != nil {
		log.Fatal(err)
	}
}
