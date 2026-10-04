package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
		respondWithError(w, http.StatusBadRequest, msgInvalidRequestBody, fmt.Errorf("decode request body: %w", err))
		return
	}

	taskID, err := uuid.Parse(r.PathValue("taskID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task ID", fmt.Errorf("parse task ID: %w", err))
		return
	}

	if taskID == req.DependencyID {
		respondWithError(w, http.StatusBadRequest, "task cannot be dependent on itself", errors.New("task id identical to dependency id"))
		return
	}

	dbDependency, err := cfg.dbQueries.GetTaskByID(r.Context(), req.DependencyID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		msg := msgServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
			msg = "dependency task does not exist"
		}
		respondWithError(w, statusCode, msg, fmt.Errorf("get dependency task by id from db: %w", err))
		return
	}
	if dbDependency.CompletedAt.Valid {
		respondWithError(w, http.StatusBadRequest, "cannot make a task dependent on a completed task", errors.New("task dependent on completed task"))
		return
	}

	dbTask, err := cfg.dbQueries.GetTaskByID(r.Context(), taskID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		msg := msgServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
			msg = "task does not exist"
		}
		respondWithError(w, statusCode, msg, fmt.Errorf("get task by id from db: %w", err))
		return
	}
	if dbTask.CompletedAt.Valid {
		respondWithError(w, http.StatusBadRequest, "cannot add a dependency to a completed task", errors.New("task already completed"))
		return
	}
	if dbDependency.DomainID != dbTask.DomainID {
		respondWithError(w, http.StatusBadRequest, "a task and its dependency must be in same domain", errors.New("task and dependency not in same domain"))
		return
	}

	//Checking if dependency leads to a cycle in the dependency graph

	visitedIDs := make(map[uuid.UUID]struct{})

	cycle, err := cfg.checkForCycle(r, visitedIDs, taskID, req.DependencyID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msgServerError, fmt.Errorf("check for cycle: %w", err))
		return
	}
	if cycle {
		respondWithError(w, http.StatusConflict, "leading to cylical dependency", errors.New("cycle detected!"))
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
		respondWithError(w, http.StatusInternalServerError, msgServerError, fmt.Errorf("create task-dependency row in db: %w", err))
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
