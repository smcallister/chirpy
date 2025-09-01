package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	// Create a random 32-byte token.
	token := make([]byte, 32)
	rand.Read(token)

	// Convert the token to a string.
	return hex.EncodeToString(token), nil
}