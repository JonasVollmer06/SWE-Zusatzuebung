package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(fussballerRouter ...http.Handler) http.Handler {
	router := chi.NewRouter()

	router.Get("/health", healthHandler)
	if len(fussballerRouter) > 0 && fussballerRouter[0] != nil {
		router.Mount("/fussballer", fussballerRouter[0])
	}

	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
