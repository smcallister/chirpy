package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/smcallister/chirpy/internal/auth"
	"github.com/smcallister/chirpy/internal/database"
	"github.com/smcallister/chirpy/internal/model"
)

func (cfg *Config) PolkaWebhookHandler(res http.ResponseWriter, req *http.Request) {
	// Add headers.
	res.Header().Add("Content-Type", "application/json")

	// Get the API key from the header.
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Validate the API key.
	if apiKey != cfg.PolkaKey {
		res.WriteHeader(401)
		res.Write([]byte("{\"error\": \"Missing or invalid token\"}"))
		return
	}

	// Decode the input.
	decoder := json.NewDecoder(req.Body)
    var event model.WebhookEvent
    err = decoder.Decode(&event)
    if err != nil {
		res.WriteHeader(400)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	// Determine the event type.
	switch event.Event {
		case "user.upgraded":
			// Upgrade the user to Red.
			cfg.handleUserUpgraded(res, req, &event.Data)
			return

		default:
			// Unknown event type.
			res.WriteHeader(204)
			return
	}
}

func (cfg *Config) handleUserUpgraded(res http.ResponseWriter, req *http.Request, data *model.WebhookEventData) {
	// Upgrade the user in the database.
	params := database.UpgradeUserToRedParams{
		ID:     	data.UserID,
		UpdatedAt: 	time.Now() }	
	
	_, err := cfg.DB.UpgradeUserToRed(req.Context(), params)
	if err != nil {
		res.WriteHeader(404)
		res.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
		return
	}

	// Write the response.
	res.WriteHeader(204)
}
