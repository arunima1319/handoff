package main

import (
	"log"
	"net/http"

	"github.com/arunima1319/handoff/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) authenticateRequest(r *http.Request) (uuid.UUID, int, string, error) {

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Could not get access token: %s", err)
		return uuid.Nil, http.StatusUnauthorized, "missing authentication credentials", err
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Could not validate access token: %s", err)
		return uuid.Nil, http.StatusUnauthorized, "invalid authentication credentials", err
	}

	return userID, 0, "", nil
}
