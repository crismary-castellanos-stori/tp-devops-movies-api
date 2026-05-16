package handler

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metricsHandler = promhttp.Handler()

func Metrics(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("METRICS_USERNAME")
	password := os.Getenv("METRICS_PASSWORD")

	if username == "" || password == "" {
		writeError(w, http.StatusServiceUnavailable, "metrics authentication is not configured")
		return
	}

	requestUsername, requestPassword, ok := r.BasicAuth()
	if !ok || requestUsername != username || requestPassword != password {
		w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	metricsHandler.ServeHTTP(w, r)
}
