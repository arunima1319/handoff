package main

import (
	"encoding/json"
	"fmt"
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
	Name string `json:"name"`
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
		return dbDomain, fmt.Errorf("begin domain owner transaction: %w", err)
	}

	defer tx.Rollback()

	qtx := cfg.dbQueries.WithTx(tx)
	dbDomain, err = qtx.CreateDomain(r.Context(), domainData)
	if err != nil {
		return dbDomain, fmt.Errorf("create domain in db: %w", err)
	}
	err = qtx.AddUserToDomain(
		r.Context(),
		database.AddUserToDomainParams{
			DomainID: dbDomain.ID,
			UserID:   dbDomain.Owner,
		})
	if err != nil {
		return dbDomain, fmt.Errorf("add owner to domain in db: %w", err)
	}

	return dbDomain, tx.Commit()
}

func (cfg *apiConfig) handlerCreateDomain(w http.ResponseWriter, r *http.Request) {

	//authenticating user

	userID, err := cfg.authenticateRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, msgFailedAuthentication, fmt.Errorf("authenticate req: %w", err))
		return
	}

	// decode request data
	req := createDomainRequest{}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, msgInvalidRequestBody, fmt.Errorf("deocde request body: %w", err))
		return
	}

	//create domain in database - will need authentication

	dbDomain, err := cfg.domainOwnerTransaction(
		r,
		database.CreateDomainParams{
			Owner: userID,
			Name:  req.Name,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msgServerError, fmt.Errorf("domain owner transaction: %w", err))
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

func (cfg *apiConfig) handlerGetUserDomains(w http.ResponseWriter, r *http.Request) {

	userID, err := cfg.authenticateRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, msgFailedAuthentication, fmt.Errorf("authenticate req: %w", err))
		return
	}

	dbDomains, err := cfg.dbQueries.GetUserDomains(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msgServerError, fmt.Errorf("get domains of user from db: %w", err))
		return
	}

	domains := []apiDomain{}
	for _, dbDomain := range dbDomains {
		domain := apiDomain{
			ID:        dbDomain.ID,
			CreatedAt: dbDomain.CreatedAt,
			UpdatedAt: dbDomain.UpdatedAt,
			Owner:     dbDomain.Owner,
			Name:      dbDomain.Name,
		}

		domains = append(domains, domain)
	}

	respondWithJSON(w, http.StatusOK, domains)
}
