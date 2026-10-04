package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arunima1319/handoff/internal/auth"
	"github.com/google/uuid"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string    `json:"access_token"`
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	req := loginReq{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, msgInvalidRequestBody, fmt.Errorf("decode request body: %w", err))
		return
	}

	payload, code, msg, err := cfg.helperLogin(r.Context(), req.Email, req.Password)
	if err != nil {
		respondWithError(w, code, msg, err)
		return
	}

	respondWithJSON(w, code, payload)
}

func (cfg *apiConfig) helperLogin(ctx context.Context, email, password string) (loginResponse, int, string, error) {

	response := loginResponse{}

	dbUser, err := cfg.dbQueries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response, http.StatusUnauthorized, msgLoginError, fmt.Errorf("email does not exist in db: %w", err)
		}
		return response, http.StatusUnauthorized, msgLoginError, fmt.Errorf("unknown error getting user by email: %w", err)
	}

	match, err := auth.CheckPasswordHash(dbUser.HashedPassword, password)
	if err != nil {
		return response, http.StatusUnauthorized, msgLoginError, fmt.Errorf("password hash: %w", err)
	}
	if match != true {
		return response, http.StatusUnauthorized, msgLoginError, errors.New("password does not match")
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, auth.JWTExpiration)
	if err != nil {
		return response, http.StatusInternalServerError, msgServerError, fmt.Errorf("creating jwt: %w", err)
	}

	response.AccessToken = accessToken
	response.ID = dbUser.ID
	response.DisplayName = dbUser.DisplayName
	response.Email = dbUser.Email

	return response, http.StatusOK, "Logged in", nil

}
