package main

import (
	"database/sql"
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

func (cfg *apiConfig) handlerCreateTaskDependency(w http.ResponseWriter, r *http.Request) {

	req := createDependencyRequest{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		statusCode, msg := reqJSONError(err)
		respondWithError(w, statusCode, msg, err)
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
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
		}
		respondWithError(w, statusCode, "Could not get dependency from database", err)
		return
	}
	if dbDependency.CompletedAt.Valid {
		respondWithError(w, http.StatusBadRequest, "Cannot make a task dependent on a completed task: ", errors.New("dependency is already completed"))
		return
	}

	dbTask, err := cfg.dbQueries.GetTaskByID(r.Context(), taskID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
		}
		respondWithError(w, statusCode, "Could not get task from database", err)
		return
	}
	if dbTask.CompletedAt.Valid {
		respondWithError(w, http.StatusBadRequest, "Cannot add a dependency to a completed task: ", errors.New("task is already completed"))
		return
	}
	if dbDependency.DomainID != dbTask.DomainID {
		respondWithError(w, http.StatusBadRequest, "A task and its dependency must be in same domain: ", errors.New("dependency not in same domain"))
		return
	}

	//Checking if dependency leads to a cycle in the dependency graph

	visitedIDs := make(map[uuid.UUID]struct{})

	cfg.checkForCycle(r, visitedIDs, taskID, req.DependencyID)

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

func (cfg *apiConfig) checkForCycle(r *http.Request, visitedIDs map[uuid.UUID]struct{}, ogTaskID uuid.UUID, dependencyID uuid.UUID) (bool, error) {

	//To check for cycle we are doing a DFS using recursion to find the original Task ID
	//If task ID is found, there is a cycle, if not, there is not

	dbDependencies, err := cfg.dbQueries.GetTaskDependenciesByTask(
		r.Context(),
		dependencyID,
	)
	if err != nil {
		return false, err
	}

	for _, row := range dbDependencies {
		_, ok := visitedIDs[row.DependencyID]
		if ok {
			continue
		}

		dependency, err := cfg.dbQueries.GetTaskByID(r.Context(), row.DependencyID)
		if err != nil {
			return false, err
		}

		if dependency.CompletedAt.Valid {
			visitedIDs[dependency.ID] = struct{}{}
			continue
		}

		if dependency.ID == ogTaskID {
			return true, nil
		}

		visitedIDs[dependency.ID] = struct{}{}

		value, err := cfg.checkForCycle(r, visitedIDs, ogTaskID, dependency.ID)
		if err != nil {
			return false, err
		}
		if value {
			return value, nil
		}
	}

	return false, nil
}
