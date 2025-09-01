package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/smcallister/chirpy/internal/auth"
	"github.com/smcallister/chirpy/internal/database"
)

type Token struct {
	AccessToken    string    `json:"token"`
}

func (cfg *Config) RefreshHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Get the token from the request.
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Get the token from the database.
	row, err := cfg.DB.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Make sure the token is not expired.
	if time.Now().After(row.ExpiresAt) || row.RevokedAt.Valid {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Generate an access token for the user.
	accessToken, err := auth.MakeJWT(row.UserID, cfg.SigningKey, time.Duration(3600) * time.Second)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	token := Token{AccessToken: accessToken}

	// Write the response.
	resBody, err := json.Marshal(token)
    if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

	res.WriteHeader(200)
	res.Write(resBody)
}


func (cfg *Config) RevokeHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Get the token from the request.
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Revoke the token in the database.
	params := database.RevokeRefreshTokenParams{
		Token:     refreshToken,
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	err = cfg.DB.RevokeRefreshToken(req.Context(), params)
	if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	// Write the response.
	res.WriteHeader(204)
}
