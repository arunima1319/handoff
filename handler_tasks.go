package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

	userID, err := cfg.authenticateRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, msgFailedAuthentication, fmt.Errorf("authenticate req: %w", err))
		return
	}

	dbTasks, err := cfg.dbQueries.GetUnblockedUserTasks(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msgServerError, fmt.Errorf("get unblocked user tasks db: %w", err))
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
		respondWithError(w, http.StatusBadRequest, "invalid domain id", fmt.Errorf("parse domain id: %w", err))
		return
	}

	dbTasks, err := cfg.dbQueries.GetTasksOfDomain(
		r.Context(),
		domainID,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msgServerError, fmt.Errorf("get tasks from domain db: %w", err))
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
		respondWithError(w, http.StatusBadRequest, "invalid domain ID", fmt.Errorf("parse domain id: %w", err))
		return
	}

	//authenticating user

	userID, err := cfg.authenticateRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, msgFailedAuthentication, fmt.Errorf("authenticate req: %w", err))
		return
	}

	req := createTaskRequest{}

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, msgInvalidRequestBody, fmt.Errorf("decode req body: %w", err))
		return
	}

	task, errorCode, errorMsg, err := cfg.helperCreateTask(r.Context(), userID, domainUUID, req.AssigneeID, req.Description)
	if err != nil {
		respondWithError(w, errorCode, errorMsg, fmt.Errorf("create task: %w", err))
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (cfg *apiConfig) helperCreateTask(ctx context.Context, userID, domainID, assigneeID uuid.UUID, description string) (apiTask, int, string, error) {

	//authorizing user

	dbDomain, err := cfg.dbQueries.GetDomainByID(ctx, domainID)
	if err != nil {
		code := http.StatusInternalServerError
		msg := msgServerError
		if errors.Is(err, sql.ErrNoRows) {
			code = http.StatusNotFound
			msg = "invalid domain ID"
		}
		return apiTask{}, code, msg, fmt.Errorf("get domain by id from db: %w", err)
	}

	if dbDomain.Owner != userID {
		return apiTask{}, http.StatusForbidden, msgForbiddenError, errors.New("user is not owner of this domain")
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
		code := http.StatusInternalServerError
		msg := msgServerError
		if pqErr := pq.As(err, pqerror.ForeignKeyViolation); pqErr != nil {
			code = http.StatusBadRequest
			msg = "invalid assignee"
		}
		return apiTask{}, code, msg, fmt.Errorf("create task in db: %w", err)

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
