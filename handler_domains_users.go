package main

import (
	"encoding/json"
	"net/http"

	"github.com/arunima1319/handoff/internal/database"
	"github.com/google/uuid"
)

type addUserToDomainRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type addUserToDomainResponse struct {
	Message string `json:"message"`
}

func (cfg *apiConfig) handlerAddUserToDomain(w http.ResponseWriter, r *http.Request) {

	req := addUserToDomainRequest{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		statusCode, msg := reqJSONError(err)
		respondWithError(w, statusCode, msg, err)
		return
	}

	domainID, err := uuid.Parse(r.PathValue("domainID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse domain ID", err)
		return
	}
	err = cfg.dbQueries.AddUserToDomain(
		r.Context(),
		database.AddUserToDomainParams{
			DomainID: domainID,
			UserID:   req.UserID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create joining row in database", err)
		return
	}

	response := addUserToDomainResponse{
		Message: "The user was added to the domain",
	}
	respondWithJSON(w, http.StatusOK, response)

}
