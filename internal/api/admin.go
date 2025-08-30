package api

import (
	"fmt"
	"net/http"
)

func (cfg *Config) MetricsHandler(res http.ResponseWriter, req *http.Request) {
	// Write the response.
	res.Header().Add("Content-Type", "text/html")
	res.WriteHeader(200)
	res.Write([]byte(
		fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",
		cfg.Hits.Load())))
}

func (cfg *Config) ResetHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")

	// Reset should only be supported in dev environments.
	if cfg.Platform != "dev" {
		res.WriteHeader(403)
		return
	}

	// Reset metrics.
	cfg.Hits.Store(0)

	// Delete all users.
	err := cfg.DB.DeleteAllUsers(req.Context())
	if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	// Write the response.
	res.WriteHeader(200)
	res.Write([]byte(fmt.Sprintf("Hits: %v", cfg.Hits.Load())))
}
