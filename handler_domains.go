package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/arunima1319/handoff/internal/database"
	"github.com/google/uuid"
)

type apiDomain struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Owner     uuid.UUID `json:"owner"`
	Name      string    `json:"name"`
}

type createDomainRequest struct {
	Owner uuid.UUID `json:"owner"`
	Name  string    `json:"name"`
}

func (cfg *apiConfig) handlerCreateDomain(w http.ResponseWriter, r *http.Request) {

	// decode request data
	req := createDomainRequest{}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request payload", err)
		return
	}

	//create domain in database - will need authentication

	dbDomain, err := cfg.dbQueries.CreateDomain(
		r.Context(),
		database.CreateDomainParams{
			Owner: req.Owner,
			Name:  req.Name,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create domain in database", err)
		return
	}

	//Responding with the domain

	domain := apiDomain{
		ID:        dbDomain.ID,
		CreatedAt: dbDomain.CreatedAt,
		UpdatedAt: dbDomain.UpdatedAt,
		Owner:     dbDomain.Owner,
		Name:      dbDomain.Name,
	}

	err = cfg.dbQueries.AddUserToDomain(
		r.Context(),
		database.AddUserToDomainParams{
			DomainID: domain.ID,
			UserID:   domain.Owner,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create joining row in database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, domain)
}
