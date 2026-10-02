package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arunima1319/handoff/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
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
	AssigneeID  uuid.UUID `json:"assignee_id"`
}

func (cfg *apiConfig) handlerGetUnblockedUserTasks(w http.ResponseWriter, r *http.Request) {

	userID, code, errorMsg, err := cfg.authenticateRequest(r)
	if err != nil {
		respondWithError(w, code, errorMsg)
		return
	}

	dbTasks, err := cfg.dbQueries.GetUnblockedUserTasks(r.Context(), userID)
	if err != nil {
		log.Printf("could not get tasks from database: %s", err)
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

	domainUUID, err := uuid.Parse(r.PathValue("domainID"))
	if err != nil {
		log.Printf("could not parse domain ID from path: %s", err)
		respondWithError(w, http.StatusInternalServerError, "could not parse domain ID")
		return
	}

	//authenticating user

	userID, errorCode, errorMsg, err := cfg.authenticateRequest(r)
	if err != nil {
		respondWithError(w, errorCode, errorMsg)
		return
	}

	req := createTaskRequest{}

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&req)
	if err != nil {
		statusCode, msg := reqJSONError(err)
		respondWithError(w, statusCode, msg)
		return
	}

	task, errorCode, errorMsg, err := cfg.helperCreateTask(r.Context(), userID, domainUUID, req.AssigneeID, req.Description)
	if err != nil {
		respondWithError(w, errorCode, errorMsg)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (cfg *apiConfig) helperCreateTask(ctx context.Context, userID, domainID, assigneeID uuid.UUID, description string) (apiTask, int, string, error) {

	//authorizing user

	dbDomain, err := cfg.dbQueries.GetDomainByID(ctx, domainID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("domain does not exist: %s", err)
			return apiTask{}, http.StatusNotFound, "invalid domain", err
		}
		log.Printf("unknown error in getting domain: %s", err)
		return apiTask{}, http.StatusInternalServerError, msgServerError, err
	}

	if dbDomain.Owner != userID {
		log.Printf("user not authorized to create task in this domain")
		return apiTask{}, http.StatusForbidden, "user forbidden to create task", err
	}

	//creating task in database

	dbTask, err := cfg.dbQueries.CreateTask(
		ctx,
		database.CreateTaskParams{
			Description: description,
			DomainID:    domainID,
			AssigneeID:  assigneeID,
		})
	if err != nil {
		if pqErr := pq.As(err, pqerror.ForeignKeyViolation); pqErr != nil {
			log.Printf("Foreign key violation in task creation: %s", err)
			return apiTask{}, http.StatusBadRequest, "invalid assignee", err
		}
		log.Printf("unknown error in creating task in database: %s", err)
		return apiTask{}, http.StatusInternalServerError, msgServerError, err

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

	return task, 0, "", nil
}
