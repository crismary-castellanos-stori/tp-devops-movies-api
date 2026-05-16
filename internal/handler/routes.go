package handler

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, movieHandler *MovieHandler) {
	r.Get("/health", Health)
	r.HandleFunc("/metrics", Metrics)
	r.Mount("/movies", movieHandler.Routes())
}
