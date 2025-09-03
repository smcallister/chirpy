package auth

import (
	"fmt"
	"net/http"
)


func GetBearerToken(headers http.Header) (string, error) {
	// Get the Authorization header.
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing Authorization header")
	}

	// Check if it is a Bearer token.
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", fmt.Errorf("Invalid Authorization header")
	}

	// Return the token.
	return authHeader[7:], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	// Get the Authorization header.
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing Authorization header")
	}

	// Check if it is an API key token.
	if len(authHeader) < 7 || authHeader[:7] != "ApiKey " {
		return "", fmt.Errorf("Invalid Authorization header")
	}

	// Return the key.
	return authHeader[7:], nil
}