package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerResetDatabase(w http.ResponseWriter, r *http.Request) {

	msg := "Could not reset database"

	tx, err := cfg.db.Begin()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	defer tx.Rollback()

	qtx := cfg.dbQueries.WithTx(tx)

	err = qtx.DeleteAllTasks(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	err = qtx.DeleteAllUsersFromDomains(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	err = qtx.DeleteAllDomains(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	err = qtx.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	respondWithJSON(w, http.StatusOK, "The database has been reset")

}
