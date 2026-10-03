package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

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

var errCouldNotLogin = errors.New("could not login")

const msgLoginError = "invalid user email or password"

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	req := loginReq{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		code, msg := reqJSONError(err)
		respondWithError(w, code, msg)
		return
	}

	payload, code, msg, err := cfg.helperLogin(r.Context(), req.Email, req.Password)
	if err != nil {
		respondWithError(w, code, msg)
		return
	}

	respondWithJSON(w, code, payload)
}

func (cfg *apiConfig) helperLogin(ctx context.Context, email, password string) (loginResponse, int, string, error) {

	response := loginResponse{}

	dbUser, err := cfg.dbQueries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Email not found in database: %v", err)
		}
		return response, http.StatusUnauthorized, msgLoginError, errCouldNotLogin
	}

	match, err := auth.CheckPasswordHash(dbUser.HashedPassword, password)
	if err != nil {
		log.Printf("Error: %v", err)
		return response, http.StatusUnauthorized, msgLoginError, errCouldNotLogin
	}
	if match != true {
		log.Printf("Password does not match")
		return response, http.StatusUnauthorized, msgLoginError, errCouldNotLogin
	}

	expirationTime, err := time.ParseDuration("1h")
	if err != nil {
		log.Printf("could not parse expiration time: %s", err)
		return response, http.StatusInternalServerError, msgServerError, errCouldNotLogin
	}
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		log.Printf("could not create access token: %s", err)
		return response, http.StatusInternalServerError, msgServerError, errCouldNotLogin
	}

	response.AccessToken = accessToken
	response.ID = dbUser.ID
	response.DisplayName = dbUser.DisplayName
	response.Email = dbUser.Email

	return response, http.StatusOK, "Logged in", nil

}
