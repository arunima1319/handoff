package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arunima1319/handoff/internal/auth"
	"github.com/arunima1319/handoff/internal/database"
	"github.com/google/uuid"
)

type apiTask struct {
	ID          uuid.UUID    `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Description string       `json:"description"`
	DomainID    uuid.UUID    `json:"domain_id"`
	AssigneeID  uuid.UUID    `json:"assignee_id"`
	CompletedAt sql.NullTime `json:"completed_at"`
}

type createTaskRequest struct {
	Description string    `json:"description"`
	DomainID    uuid.UUID `json:"domain_id"`
	AssigneeID  uuid.UUID `json:"assignee_id"`
}

func (cfg *apiConfig) handlerGetUnblockedUserTasks(w http.ResponseWriter, r *http.Request) {

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Could not get access token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "missing authentication credentials")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Could not validate access token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "invalid authentication credentials")
	}

	dbTasks, err := cfg.dbQueries.GetUnblockedUserTasks(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get tasks")
		return
	}

	tasks := []apiTask{}
	for _, dbTask := range dbTasks {
		task := apiTask{
			ID:          dbTask.ID,
			CreatedAt:   dbTask.CreatedAt,
			UpdatedAt:   dbTask.UpdatedAt,
			Description: dbTask.Description,
			DomainID:    dbTask.DomainID,
			AssigneeID:  dbTask.AssigneeID,
			CompletedAt: dbTask.CompletedAt,
		}

		tasks = append(tasks, task)
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (cfg *apiConfig) handlerGetTasksOfDomain(w http.ResponseWriter, r *http.Request) {

	domainID, err := uuid.Parse(r.PathValue("domainID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse domain ID")
		return
	}

	dbTasks, err := cfg.dbQueries.GetTasksOfDomain(
		r.Context(),
		domainID,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get tasks from domain")
		return
	}

	tasks := []apiTask{}
	for _, dbTask := range dbTasks {
		task := apiTask{
			ID:          dbTask.ID,
			CreatedAt:   dbTask.CreatedAt,
			UpdatedAt:   dbTask.UpdatedAt,
			Description: dbTask.Description,
			DomainID:    dbTask.DomainID,
			AssigneeID:  dbTask.AssigneeID,
			CompletedAt: dbTask.CompletedAt,
		}
		tasks = append(tasks, task)

	}

	respondWithJSON(w, http.StatusOK, tasks)

}

func (cfg *apiConfig) handlerCreateTask(w http.ResponseWriter, r *http.Request) {

	req := createTaskRequest{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		statusCode, msg := reqJSONError(err)
		respondWithError(w, statusCode, msg)
		return
	}

	dbTask, err := cfg.dbQueries.CreateTask(
		r.Context(),
		database.CreateTaskParams{
			Description: req.Description,
			DomainID:    req.DomainID,
			AssigneeID:  req.AssigneeID,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create task in database")
		return
	}

	task := apiTask{
		ID:          dbTask.ID,
		CreatedAt:   dbTask.CreatedAt,
		UpdatedAt:   dbTask.UpdatedAt,
		Description: dbTask.Description,
		DomainID:    dbTask.DomainID,
		AssigneeID:  dbTask.AssigneeID,
		CompletedAt: dbTask.CompletedAt,
	}

	respondWithJSON(w, http.StatusOK, task)
}
