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

func (cfg *apiConfig) domainOwnerTransaction(r *http.Request, domainData database.CreateDomainParams) (database.Domain, error) {

	/*
		Helper function to ensure that creating a domain
		and adding the owner to the domain as a user is
		a transaction
	*/

	dbDomain := database.Domain{}

	tx, err := cfg.db.Begin()
	if err != nil {
		return dbDomain, err
	}

	defer tx.Rollback()

	qtx := cfg.dbQueries.WithTx(tx)
	dbDomain, err = qtx.CreateDomain(r.Context(), domainData)
	if err != nil {
		return dbDomain, err
	}
	err = qtx.AddUserToDomain(
		r.Context(),
		database.AddUserToDomainParams{
			DomainID: dbDomain.ID,
			UserID:   dbDomain.Owner,
		})
	if err != nil {
		return dbDomain, err
	}

	return dbDomain, tx.Commit()
}

func (cfg *apiConfig) handlerCreateDomain(w http.ResponseWriter, r *http.Request) {

	// decode request data
	req := createDomainRequest{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		statusCode, msg := reqJSONError(err)
		respondWithError(w, statusCode, msg, err)
		return
	}

	//create domain in database - will need authentication

	dbDomain, err := cfg.domainOwnerTransaction(
		r,
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

	respondWithJSON(w, http.StatusOK, domain)
}
