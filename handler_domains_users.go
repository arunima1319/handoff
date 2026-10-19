package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arunima1319/handoff/internal/database"
	"github.com/google/uuid"
)

type addUserToDomainRequest struct {
	Email string `json:"email"`
}

type addUserToDomainResponse struct {
	Message string `json:"message"`
}

func (cfg *apiConfig) handlerAddUserToDomain(w http.ResponseWriter, r *http.Request) {

	req := addUserToDomainRequest{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, msgInvalidRequestBody, fmt.Errorf("decode request body: %w", err))
		return
	}

	domainID, err := uuid.Parse(r.PathValue("domainID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid domain ID", fmt.Errorf("parse domain id: %w", err))
		return
	}

	dbDomain, err := cfg.dbQueries.GetDomainByID(r.Context(), domainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get domain", fmt.Errorf("get db domain from id: %w", err))
	}
	userID, err := cfg.authenticateRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, msgFailedAuthentication, fmt.Errorf("authenticate req: %w", err))
	}
	if userID != dbDomain.Owner {
		respondWithError(w, http.StatusForbidden, msgForbiddenError, errors.New("non-owner attempt to add user to domain"))
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msgServerError, fmt.Errorf("get user from email: %w", err))
	}

	err = cfg.dbQueries.AddUserToDomain(
		r.Context(),
		database.AddUserToDomainParams{
			DomainID: domainID,
			UserID:   dbUser.ID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msgServerError, fmt.Errorf("add user to domain in database: %w", err))
		return
	}

	response := addUserToDomainResponse{
		Message: "The user was added to the domain",
	}
	respondWithJSON(w, http.StatusOK, response)

}
