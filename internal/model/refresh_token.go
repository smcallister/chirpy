package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/smcallister/chirpy/internal/auth"
	"github.com/smcallister/chirpy/internal/database"
)

func NewRefreshToken(context context.Context, userID uuid.UUID, db *database.Queries) (*database.RefreshToken, error) {
	// Generate a random refresh token.
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return nil, err
	}
	
	// Add the refresh token to the database.
	currentTime := time.Now()
	params := database.CreateRefreshTokenParams{
		refreshToken,
		currentTime,
		currentTime,
		userID,
		currentTime.Add(24 * 60 * time.Hour)}

	row, err := db.CreateRefreshToken(context, params)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

