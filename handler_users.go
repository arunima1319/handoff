package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arunima1319/handoff/internal/auth"
	"github.com/arunima1319/handoff/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

type apiUser struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
}

type createUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (cfg *apiConfig) handlerGetUsersOfDomain(w http.ResponseWriter, r *http.Request) {

	domainID, err := uuid.Parse(r.PathValue("domainID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse domain ID")
		return
	}

	dbUsers, err := cfg.dbQueries.GetUsersOfDomain(r.Context(), domainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get users from database")
		return
	}

	users := []apiUser{}

	for _, dbUser := range dbUsers {
		user := apiUser{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			DisplayName: dbUser.DisplayName,
		}

		users = append(users, user)
	}

	respondWithJSON(w, http.StatusOK, users)
}
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	req := createUserRequest{}

	// Deserializing the request payload
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		statusCode, msg := reqJSONError(err)
		respondWithError(w, statusCode, msg)
		return
	}

	user, code, msg, err := cfg.helperCreateUser(r.Context(), req)
	if err != nil {
		respondWithError(w, code, msg)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) helperCreateUser(ctx context.Context, req createUserRequest) (apiUser, int, string, error) {

	//Hashing password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Could not hash password: %s", err)
		return apiUser{}, http.StatusInternalServerError, "internal server error", err
	}

	// Creating the user in the database
	dbUser, err := cfg.dbQueries.CreateUser(
		ctx,
		database.CreateUserParams{
			Email:          req.Email,
			DisplayName:    req.DisplayName,
			HashedPassword: hashedPassword,
		})
	if err != nil {
		if pqErr := pq.As(err, pqerror.UniqueViolation); pqErr != nil {
			log.Printf("email %s already exists", req.Email)
			return apiUser{}, http.StatusConflict, "email already exists", err
		}

		log.Printf("unknown error in creating user in database: %s", err)
		return apiUser{}, http.StatusInternalServerError, msgServerError, err
	}

	//Writing the response

	user := apiUser{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		DisplayName: dbUser.DisplayName,
	}

	return user, http.StatusOK, "", nil
}
