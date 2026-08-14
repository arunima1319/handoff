package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arunima1319/handoff/internal/database"
	"github.com/google/uuid"
)

type createDependencyRequest struct {
	DependencyID uuid.UUID `json:"dependency_id"`
}

type createDependencyResponse struct {
	Message string `json:"message"`
}

type apiTaskDependency struct {
	taskID       uuid.UUID
	dependencyID uuid.UUID
}

func (cfg *apiConfig) handlerCreateTaskDependency(w http.ResponseWriter, r *http.Request) {

	req := createDependencyRequest{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode request body", err)
		return
	}

	taskID, err := uuid.Parse(r.PathValue("taskID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not parse Task ID", err)
		return
	}

	if taskID == req.DependencyID {
		respondWithError(w, http.StatusBadRequest, "Task cannot be dependent on itself", errors.New("Task and Dependency have same ID"))
		return
	}

	dbDependency, err := cfg.dbQueries.GetTaskByID(r.Context(), req.DependencyID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get dependency from database", err)
		return
	}
	if dbDependency.CompletedAt.Valid {
		respondWithError(w, http.StatusBadRequest, "Cannot make a task dependent on a completed task", errors.New("Dependency is already completed"))
		return
	}

	dbTask, err := cfg.dbQueries.GetTaskByID(r.Context(), taskID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get task from database", err)
		return
	}
	if dbTask.CompletedAt.Valid {
		respondWithError(w, http.StatusBadRequest, "Cannot add a dependency to a completed task", errors.New("Task is already completed"))
		return
	}
	if dbDependency.DomainID != dbTask.DomainID {
		respondWithError(w, http.StatusBadRequest, "A task and its dependency must be in same domain", errors.New("Dependency not in same domain"))
		return
	}

	// Creating Task Dependency in Database
	err = cfg.dbQueries.CreateTaskDependency(
		r.Context(),
		database.CreateTaskDependencyParams{
			TaskID:       taskID,
			DependencyID: req.DependencyID,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not add dependency in database", err)
		return
	}

	msg := createDependencyResponse{
		Message: "Dependency has been created",
	}

	respondWithJSON(w, http.StatusOK, msg)

}
