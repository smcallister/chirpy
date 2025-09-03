package api

import (
	"net/http"
	"sync/atomic"

	"github.com/smcallister/chirpy/internal/database"
)

type Config struct {
	Hits 		atomic.Int32
	DB 			*database.Queries
	SigningKey 	string
	Platform	string
	PolkaKey    string
}

func (cfg *Config) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.Hits.Add(1)
		next.ServeHTTP(res, req)
	})
}
