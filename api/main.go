package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/client"
	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/config"
	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/handler"
	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/monitoring"
	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	tmdbClient := client.NewTMDBClient(cfg.TMDBBaseURL, cfg.TMDBApiKey)
	movieService := service.NewMovieService(tmdbClient)
	movieHandler := handler.NewMovieHandler(movieService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(monitoring.Middleware)

	handler.RegisterRoutes(r, movieHandler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrCh := make(chan error, 1)

	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("Received signal %s, shutting down server...", sig)
	case err := <-serverErrCh:
		log.Printf("Server error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
