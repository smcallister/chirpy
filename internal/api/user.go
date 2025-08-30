package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
