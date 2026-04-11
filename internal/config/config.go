package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	TMDBApiKey  string
	TMDBBaseURL string
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEY environment variable is required")
	}

	baseURL := os.Getenv("TMDB_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.themoviedb.org/3"
	}

	return &Config{
		Port:        port,
		TMDBApiKey:  apiKey,
		TMDBBaseURL: baseURL,
	}, nil
}
