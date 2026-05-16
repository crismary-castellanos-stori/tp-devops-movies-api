package handler

import "net/http"

// DebugError provides a predictable 500 response for monitoring demos.
func DebugError(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusInternalServerError, "intentional debug error")
}
