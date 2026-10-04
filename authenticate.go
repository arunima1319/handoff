package main

import (
	"fmt"
	"net/http"

	"github.com/arunima1319/handoff/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) authenticateRequest(r *http.Request) (uuid.UUID, error) {

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get bearer token: %w", err)
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("validate jwt: %w", err)
	}

	return userID, nil
}
