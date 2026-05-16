package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestMetricsRequiresAuthentication(t *testing.T) {
	t.Setenv("METRICS_USERNAME", "metrics-user")
	t.Setenv("METRICS_PASSWORD", "metrics-pass")

	r := chi.NewRouter()
	r.HandleFunc("/metrics", Metrics)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestMetricsReturnsMetricsWithValidCredentials(t *testing.T) {
	t.Setenv("METRICS_USERNAME", "metrics-user")
	t.Setenv("METRICS_PASSWORD", "metrics-pass")

	r := chi.NewRouter()
	r.HandleFunc("/metrics", Metrics)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	req.SetBasicAuth("metrics-user", "metrics-pass")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected metrics response body")
	}
}

func TestMetricsReturnsServiceUnavailableWhenCredentialsAreMissing(t *testing.T) {
	os.Unsetenv("METRICS_USERNAME")
	os.Unsetenv("METRICS_PASSWORD")

	r := chi.NewRouter()
	r.HandleFunc("/metrics", Metrics)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}
}
