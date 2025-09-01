package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var signingKey = "my_secret_key"

func TestValidJWT(t *testing.T) {
	// Generate the user ID.
	id := uuid.New()
	
	// Generate the JWT.
	token, err := MakeJWT(id, signingKey, time.Minute * 5)
	if err != nil {
		t.Fatalf("MakeJWT returned an unexpected error: %v", err)
	}

	// Validate the JWT.
	idFromToken, err := ValidateJWT(token, signingKey)
	if err != nil {
		t.Fatalf("ValidateJWT returned an unexpected error: %v", err)
	}

	if id != idFromToken {
		t.Fatalf("Expected %v but got %v from ValidateJWT", id, idFromToken)
	}
}

func TestExpiredJWT(t *testing.T) {
	// Generate the user ID.
	id := uuid.New()
	
	// Generate the JWT.
	token, err := MakeJWT(id, signingKey, time.Second * 1)
	if err != nil {
		t.Fatalf("MakeJWT returned an unexpected error: %v", err)
	}

	// Wait for the JWT to expire.
	time.Sleep(time.Second * 2)

	// Validate the JWT and make sure it fails.
	_, err 	= ValidateJWT(token, signingKey)
	if err == nil {
		t.Fatalf("ValidateJWT returned success, but expected an error")
	}
}

func TestJWTWithInvalidSigningKey(t *testing.T) {
	// Generate the user ID.
	id := uuid.New()
	
	// Generate the JWT.
	token, err := MakeJWT(id, signingKey, time.Minute * 5)
	if err != nil {
		t.Fatalf("MakeJWT returned an unexpected error: %v", err)
	}

	// Validate the JWT and make sure it fails.
	_, err = ValidateJWT(token, signingKey + "_bad")
	if err == nil {
		t.Fatalf("ValidateJWT returned success, but expected an error")
	}
}
