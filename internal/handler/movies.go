package handler

import (
	"errors"
	"net/http"

	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type MovieHandler struct {
	movieService *service.MovieService
}

func NewMovieHandler(movieService *service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (h *MovieHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/search", h.Search)
	r.Get("/{id}", h.GetByID)

	return r
}

func (h *MovieHandler) Search(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	result, err := h.movieService.Search(r.Context(), title)
	if err != nil {
		if errors.Is(err, service.ErrEmptyTitle) {
			writeError(w, http.StatusBadRequest, "title query parameter is required")
			return
		}
		handleMovieError(w, err, "failed to search movies")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	movie, err := h.movieService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			writeError(w, http.StatusBadRequest, "invalid movie id")
			return
		}
		if errors.Is(err, service.ErrMovieNotFound) {
			writeError(w, http.StatusNotFound, "movie not found")
			return
		}
		handleMovieError(w, err, "failed to get movie")
		return
	}

	writeJSON(w, http.StatusOK, movie)
}

func handleMovieError(w http.ResponseWriter, err error, fallbackMessage string) {
	var upstreamErr *service.UpstreamError
	if errors.As(err, &upstreamErr) {
		if upstreamErr.RetryAfter != "" {
			w.Header().Set("Retry-After", upstreamErr.RetryAfter)
		}

		status, message := mapUpstreamError(upstreamErr, fallbackMessage)
		writeError(w, status, message)
		return
	}

	writeError(w, http.StatusBadGateway, fallbackMessage)
}

func mapUpstreamError(err *service.UpstreamError, fallbackMessage string) (int, string) {
	switch err.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return http.StatusBadGateway, "upstream movie service rejected the request"
	case http.StatusTooManyRequests:
		return http.StatusTooManyRequests, "upstream movie service rate limit exceeded"
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return http.StatusBadGateway, "upstream movie service is unavailable"
	default:
		if err.StatusCode >= 500 {
			return http.StatusBadGateway, "upstream movie service is unavailable"
		}
		return http.StatusBadGateway, fallbackMessage
	}
}
