package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/smcallister/chirpy/internal/api"
	"github.com/smcallister/chirpy/internal/database"
)

import _ "github.com/lib/pq"

func healthzHandler(res http.ResponseWriter, req *http.Request) {
	// Write the response.
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func main() {
	// Initialization.
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Failed to open DB %s: %v", dbURL, err)
	}

	apiCfg := api.Config{
		DB: database.New(db),
		Platform: os.Getenv("PLATFORM"),
		SigningKey: os.Getenv("JWT_SIGNING_KEY")}

	// Set up the handlers.
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("/assets/", http.FileServer(http.Dir(".")))

	mux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.MetricsHandler))
	mux.Handle("POST /admin/reset", http.HandlerFunc(apiCfg.ResetHandler))
	
	mux.HandleFunc("GET /api/healthz", http.HandlerFunc(healthzHandler))
	mux.HandleFunc("POST /api/users", http.HandlerFunc(apiCfg.CreateUserHandler))
	mux.HandleFunc("PUT /api/users", http.HandlerFunc(apiCfg.UpdateUserHandler))
	
	mux.HandleFunc("POST /api/login", http.HandlerFunc(apiCfg.LoginHandler))
	mux.HandleFunc("POST /api/refresh", http.HandlerFunc(apiCfg.RefreshHandler))
	mux.HandleFunc("POST /api/revoke", http.HandlerFunc(apiCfg.RevokeHandler))

	mux.HandleFunc("GET /api/chirps", http.HandlerFunc(apiCfg.GetChirpsHandler))
	mux.HandleFunc("GET /api/chirps/{id}", http.HandlerFunc(apiCfg.GetChirpHandler))
	mux.HandleFunc("POST /api/chirps", http.HandlerFunc(apiCfg.CreateChirpHandler))
	mux.HandleFunc("DELETE /api/chirps/{id}", http.HandlerFunc(apiCfg.DeleteChirpHandler))

	// Create the server.	
	server := http.Server{Handler: mux, Addr: ":8080"}

	// Listen for requests.
	server.ListenAndServe()
}