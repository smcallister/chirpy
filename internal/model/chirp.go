package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/smcallister/chirpy/internal/database"
)

func NewChirp(context context.Context, r io.Reader, db *database.Queries) (*database.Chirp, error) {
	// Decode the chirp.
	decoder := json.NewDecoder(r)
    var chirp database.Chirp
    err := decoder.Decode(&chirp)
    if err != nil {
		return nil, err
	}

	// Validate the chirp.
	if len(chirp.Body) > 140 {
		return nil, fmt.Errorf("Chirp is too long")
	}

	// Replace profanity.
	words := strings.Split(chirp.Body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			words[i] = "****"
		}
	}

	chirp.Body = strings.Join(words, " ")

	// Create the chirp.
	currentTime := time.Now()
	params := database.CreateChirpParams{
		uuid.New(),
		currentTime,
		currentTime,
		chirp.Body,
		chirp.UserID }

	row, err := db.CreateChirp(context, params)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

