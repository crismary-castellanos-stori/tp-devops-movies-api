package service

import (
	"context"
	"errors"
	"testing"

	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/client"
	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/model"
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

func TestGetByIDRejectsNonNumericID(t *testing.T) {
	called := false
	svc := NewMovieService(&stubMovieClient{
		getMovieFn: func(ctx context.Context, id string) (*model.TMDBMovieDetail, error) {
			called = true
			return nil, nil
		},
	})

	_, err := svc.GetByID(context.Background(), "abc")
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}

	if called {
		t.Fatal("expected client not to be called for invalid id")
	}
}

func TestGetByIDMapsNotFound(t *testing.T) {
	svc := NewMovieService(&stubMovieClient{
		getMovieFn: func(ctx context.Context, id string) (*model.TMDBMovieDetail, error) {
			return nil, client.ErrNotFound
		},
	})

	_, err := svc.GetByID(context.Background(), "123")
	if !errors.Is(err, ErrMovieNotFound) {
		t.Fatalf("expected ErrMovieNotFound, got %v", err)
	}
}

func TestSearchMapsUpstreamError(t *testing.T) {
	svc := NewMovieService(&stubMovieClient{
		searchMoviesFn: func(ctx context.Context, title string) (*model.TMDBSearchResponse, error) {
			return nil, &client.UpstreamError{
				StatusCode: 429,
				Message:    "rate limit",
				RetryAfter: "60",
			}
		},
	})

	_, err := svc.Search(context.Background(), "batman")
	if err == nil {
		t.Fatal("expected an error")
	}

	var upstreamErr *UpstreamError
	if !errors.As(err, &upstreamErr) {
		t.Fatalf("expected UpstreamError, got %T", err)
	}

	if upstreamErr.StatusCode != 429 {
		t.Fatalf("expected status 429, got %d", upstreamErr.StatusCode)
	}

	if upstreamErr.RetryAfter != "60" {
		t.Fatalf("expected retry-after 60, got %q", upstreamErr.RetryAfter)
	}
}

func TestGetByIDNormalizesNumericID(t *testing.T) {
	var receivedID string
	svc := NewMovieService(&stubMovieClient{
		getMovieFn: func(ctx context.Context, id string) (*model.TMDBMovieDetail, error) {
			receivedID = id
			return &model.TMDBMovieDetail{ID: 7, Title: "Se7en"}, nil
		},
	})

	movie, err := svc.GetByID(context.Background(), "007")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if receivedID != "7" {
		t.Fatalf("expected normalized id 7, got %q", receivedID)
	}

	if movie.ID != 7 || movie.Title != "Se7en" {
		t.Fatalf("unexpected movie result: %+v", movie)
	}
}
