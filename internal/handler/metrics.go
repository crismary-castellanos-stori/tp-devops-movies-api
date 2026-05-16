package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metricsHandler = promhttp.Handler()

func Metrics(w http.ResponseWriter, r *http.Request) {
	metricsHandler.ServeHTTP(w, r)
}
