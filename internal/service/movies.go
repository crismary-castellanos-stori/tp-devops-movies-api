package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/client"
	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/model"
)

var (
	ErrEmptyTitle    = errors.New("title parameter is required")
	ErrInvalidID     = errors.New("invalid movie id")
	ErrMovieNotFound = errors.New("movie not found")
)

type UpstreamError struct {
	StatusCode int
	Message    string
	RetryAfter string
}

func (e *UpstreamError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("upstream service returned status %d", e.StatusCode)
	}

	return fmt.Sprintf("upstream service returned status %d: %s", e.StatusCode, e.Message)
}

type movieClient interface {
	SearchMovies(ctx context.Context, title string) (*model.TMDBSearchResponse, error)
	GetMovie(ctx context.Context, id string) (*model.TMDBMovieDetail, error)
}

type MovieService struct {
	tmdbClient movieClient
}

func NewMovieService(tmdbClient movieClient) *MovieService {
	return &MovieService{
		tmdbClient: tmdbClient,
	}
}

func (s *MovieService) Search(ctx context.Context, title string) (*model.SearchResponse, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}

	result, err := s.tmdbClient.SearchMovies(ctx, title)
	if err != nil {
		var upstreamErr *client.UpstreamError
		if errors.As(err, &upstreamErr) {
			return nil, &UpstreamError{
				StatusCode: upstreamErr.StatusCode,
				Message:    upstreamErr.Message,
				RetryAfter: upstreamErr.RetryAfter,
			}
		}
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	movies := make([]model.Movie, 0, len(result.Results))
	for _, m := range result.Results {
		movies = append(movies, m.ToMovie())
	}

	return &model.SearchResponse{
		Results: movies,
		Total:   result.TotalResults,
	}, nil
}

func (s *MovieService) GetByID(ctx context.Context, id string) (*model.Movie, error) {
	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID <= 0 {
		return nil, ErrInvalidID
	}

	result, err := s.tmdbClient.GetMovie(ctx, strconv.Itoa(parsedID))
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil, ErrMovieNotFound
		}
		var upstreamErr *client.UpstreamError
		if errors.As(err, &upstreamErr) {
			return nil, &UpstreamError{
				StatusCode: upstreamErr.StatusCode,
				Message:    upstreamErr.Message,
				RetryAfter: upstreamErr.RetryAfter,
			}
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	movie := result.ToMovie()
	return &movie, nil
}
