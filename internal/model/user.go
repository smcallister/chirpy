package model

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/smcallister/chirpy/internal/database"
)

func NewUser(context context.Context, r io.Reader, db *database.Queries) (*database.User, error) {
	// Decode the user.
	decoder := json.NewDecoder(r)
    var user database.User
    err := decoder.Decode(&user)
    if err != nil {
		return nil, err
	}

	// Create the user.
	currentTime := time.Now()
	params := database.CreateUserParams{
		uuid.New(),
		currentTime,
		currentTime,
		user.Email }

	row, err := db.CreateUser(context, params)
	if err != nil {
		return nil, err
	}

	return &row, nil
}