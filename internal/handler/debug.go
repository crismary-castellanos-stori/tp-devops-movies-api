package handler

import (
	"net/http"
	"strconv"
	"time"
)

// DebugError provides a predictable 500 response for monitoring demos.
func DebugError(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusInternalServerError, "intentional debug error")
}

// DebugSlow adds a controlled delay so latency panels can show clear spikes.
func DebugSlow(w http.ResponseWriter, r *http.Request) {
	delayMS := 1200

	if rawDelay := r.URL.Query().Get("ms"); rawDelay != "" {
		parsedDelay, err := strconv.Atoi(rawDelay)
		if err != nil || parsedDelay <= 0 || parsedDelay > 10000 {
			writeError(w, http.StatusBadRequest, "ms query parameter must be between 1 and 10000")
			return
		}
		delayMS = parsedDelay
	}

	select {
	case <-time.After(time.Duration(delayMS) * time.Millisecond):
		writeJSON(w, http.StatusOK, map[string]int{"delay_ms": delayMS})
	case <-r.Context().Done():
		writeError(w, http.StatusRequestTimeout, "request canceled")
	}
}
