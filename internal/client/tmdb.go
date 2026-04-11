package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crismary-castellanos-stori/tp-devops-movies-api/internal/model"
)

var ErrNotFound = errors.New("resource not found")

type UpstreamError struct {
	StatusCode int
	Message    string
	RetryAfter string
}

func (e *UpstreamError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("tmdb api returned status %d", e.StatusCode)
	}

	return fmt.Sprintf("tmdb api returned status %d: %s", e.StatusCode, e.Message)
}

type TMDBClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewTMDBClient(baseURL, apiKey string) *TMDBClient {
	return &TMDBClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (c *TMDBClient) SearchMovies(ctx context.Context, title string) (*model.TMDBSearchResponse, error) {
	endpoint := fmt.Sprintf("%s/search/movie?query=%s", c.baseURL, url.QueryEscape(title))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, newUpstreamError(resp)
	}

	var result model.TMDBSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *TMDBClient) GetMovie(ctx context.Context, id string) (*model.TMDBMovieDetail, error) {
	endpoint := fmt.Sprintf("%s/movie/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newUpstreamError(resp)
	}

	var result model.TMDBMovieDetail
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *TMDBClient) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func newUpstreamError(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
	if err != nil {
		return &UpstreamError{
			StatusCode: resp.StatusCode,
			RetryAfter: resp.Header.Get("Retry-After"),
		}
	}

	message := strings.TrimSpace(string(body))
	if message == "" {
		return &UpstreamError{
			StatusCode: resp.StatusCode,
			RetryAfter: resp.Header.Get("Retry-After"),
		}
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err == nil {
		if statusMessage, ok := payload["status_message"].(string); ok && statusMessage != "" {
			message = statusMessage
		} else if status, ok := payload["status"].(string); ok && status != "" {
			message = status
		}
	}

	return &UpstreamError{
		StatusCode: resp.StatusCode,
		Message:    message,
		RetryAfter: resp.Header.Get("Retry-After"),
	}
}
