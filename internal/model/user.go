package model

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/smcallister/chirpy/internal/auth"
	"github.com/smcallister/chirpy/internal/database"
)

type UserInput struct {
	Email          string    `json:"email"`
	Password 	   string    `json:"password"`
}

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	AccessToken    string    `json:"token"`
	RefreshToken   string    `json:"refresh_token"`
}

func NewUser(context context.Context, r io.Reader, db *database.Queries) (*User, error) {
	// Decode the input.
	decoder := json.NewDecoder(r)
    var input UserInput
    err := decoder.Decode(&input)
    if err != nil {
		return nil, err
	}

	// Hash the password.
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create the user.
	currentTime := time.Now()
	params := database.CreateUserParams{
		uuid.New(),
		currentTime,
		currentTime,
		input.Email,
		hashedPassword }

	row, err := db.CreateUser(context, params)
	if err != nil {
		return nil, err
	}

	// Convert the row into an API object.	
	user := User{
		ID: row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Email: row.Email}

	return &user, nil
}