package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/smcallister/chirpy/internal/auth"
	"github.com/smcallister/chirpy/internal/model"
)

func (cfg *Config) CreateUserHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Create the user.
	user, err := model.NewUser(req.Context(), req.Body, cfg.DB)
	if err != nil {
		res.WriteHeader(400)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	// Write the response.
	resBody, err := json.Marshal(user)
    if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

	res.WriteHeader(201)
	res.Write(resBody)
}

func (cfg *Config) LoginHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Decode the user.
	decoder := json.NewDecoder(req.Body)
    var login model.UserInput
    err := decoder.Decode(&login)
    if err != nil {
		res.WriteHeader(400)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	// Get the user from the database.
	row, err := cfg.DB.GetUserByEmail(req.Context(), login.Email)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Incorrect email or password\"}"))
		return
	}

	user := model.User{
		ID: row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Email: row.Email,
	}

	// Make sure the passwords match.
	err = auth.CheckPasswordHash(login.Password, row.HashedPassword)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Incorrect email or password\"}"))
		return
	}

	// Generate an access token for the user.
	accessToken, err := auth.MakeJWT(user.ID, cfg.SigningKey, time.Duration(3600) * time.Second)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	user.AccessToken = accessToken

	// Generate a refresh token for the user.
	refreshToken, err := model.NewRefreshToken(req.Context(), user.ID, cfg.DB)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	user.RefreshToken = refreshToken.Token

	// Write the response.
	resBody, err := json.Marshal(user)
    if err != nil {
        res.WriteHeader(500)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
    }

	res.WriteHeader(200)
	res.Write(resBody)
}