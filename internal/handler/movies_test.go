package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/client"
	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/model"
	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type stubMovieClient struct {
	searchMoviesFn func(ctx context.Context, title string) (*model.TMDBSearchResponse, error)
	getMovieFn     func(ctx context.Context, id string) (*model.TMDBMovieDetail, error)
}

func (s *stubMovieClient) SearchMovies(ctx context.Context, title string) (*model.TMDBSearchResponse, error) {
	if s.searchMoviesFn == nil {
		return nil, nil
	}

	return s.searchMoviesFn(ctx, title)
}

func (s *stubMovieClient) GetMovie(ctx context.Context, id string) (*model.TMDBMovieDetail, error) {
	if s.getMovieFn == nil {
		return nil, nil
	}

	return s.getMovieFn(ctx, id)
}

func TestGetByIDReturnsBadRequestForInvalidID(t *testing.T) {
	h := NewMovieHandler(service.NewMovieService(&stubMovieClient{}))
	r := chi.NewRouter()
	r.Mount("/movies", h.Routes())

	req := httptest.NewRequest(http.MethodGet, "/movies/not-a-number", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	var body ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if body.Error != "invalid movie id" {
		t.Fatalf("expected invalid id message, got %q", body.Error)
	}
}

func TestSearchReturnsTooManyRequestsAndRetryAfter(t *testing.T) {
	h := NewMovieHandler(service.NewMovieService(&stubMovieClient{
		searchMoviesFn: func(ctx context.Context, title string) (*model.TMDBSearchResponse, error) {
			return nil, &client.UpstreamError{
				StatusCode: http.StatusTooManyRequests,
				Message:    "rate limit",
				RetryAfter: "30",
			}
		},
	}))
	r := chi.NewRouter()
	r.Mount("/movies", h.Routes())

	req := httptest.NewRequest(http.MethodGet, "/movies/search?title=batman", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", rec.Code)
	}

	if got := rec.Header().Get("Retry-After"); got != "30" {
		t.Fatalf("expected Retry-After 30, got %q", got)
	}

	var body ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if body.Error != "upstream movie service rate limit exceeded" {
		t.Fatalf("unexpected error message: %q", body.Error)
	}
}

func TestGetByIDReturnsMovie(t *testing.T) {
	h := NewMovieHandler(service.NewMovieService(&stubMovieClient{
		getMovieFn: func(ctx context.Context, id string) (*model.TMDBMovieDetail, error) {
			return &model.TMDBMovieDetail{
				ID:          550,
				Title:       "Fight Club",
				Overview:    "Test overview",
				ReleaseDate: "1999-10-15",
				PosterPath:  "/poster.jpg",
				VoteAverage: 8.4,
			}, nil
		},
	}))
	r := chi.NewRouter()
	r.Mount("/movies", h.Routes())

	req := httptest.NewRequest(http.MethodGet, "/movies/550", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body model.Movie
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if body.ID != 550 || body.Title != "Fight Club" {
		t.Fatalf("unexpected movie response: %+v", body)
	}
}
