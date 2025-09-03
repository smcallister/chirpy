package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/smcallister/chirpy/internal/auth"
	"github.com/smcallister/chirpy/internal/database"
	"github.com/smcallister/chirpy/internal/model"
)

func (cfg *Config) CreateChirpHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Get the token from the request.
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Validate the token.
	userID, err := auth.ValidateJWT(token, cfg.SigningKey)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Create the chirp.
	chirp, err := model.NewChirp(req.Context(), userID, req.Body, cfg.DB)
	if err != nil {
		res.WriteHeader(400)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	// Write the response.
	resBody, err := json.Marshal(chirp)
    if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

	res.WriteHeader(201)
	res.Write(resBody)
}

func (cfg *Config) GetChirpsHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Get all chirps.
	chirps, err := cfg.DB.GetAllChirps(req.Context())
    if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

	// Write the response.
	resBody, err := json.Marshal(chirps)
    if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

	res.WriteHeader(200)
	res.Write(resBody)
}

func (cfg *Config) GetChirpHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Get the requested chirp.
	id, err := uuid.Parse(req.PathValue("id"));
    if err != nil {
        res.WriteHeader(400)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

	chirp, err := cfg.DB.GetChirpByID(req.Context(), id)
    if err != nil {
        res.WriteHeader(404)
		return
	}

	// Write the response.
	resBody, err := json.Marshal(chirp)
    if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

	res.WriteHeader(200)
	res.Write(resBody)
}

func (cfg *Config) DeleteChirpHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Get the token from the request.
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Validate the token.
	userID, err := auth.ValidateJWT(token, cfg.SigningKey)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}
	
	// Get the chirp and ensure that it belongs to the authenticated user.
	id, err := uuid.Parse(req.PathValue("id"));
	chirp, err := cfg.DB.GetChirpByID(req.Context(), id)
    if err != nil {
        res.WriteHeader(404)
		return
	}

	if chirp.UserID != userID {
		res.WriteHeader(403)
		return
	}

	// Delete the chirp.
    if err != nil {
        res.WriteHeader(400)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

 	params := database.DeleteChirpParams{
		ID:     id,
		UserID: userID}

	err = cfg.DB.DeleteChirp(req.Context(), params)
    if err != nil {
        res.WriteHeader(500)
		return
	}

	// Write the response.
	res.WriteHeader(204)
}
