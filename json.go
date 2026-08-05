package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload any) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error in Marshaling response payload: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Error in writing data to connection: %s", err)
	}

}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {

	type errorResponse struct {
		Error string `json:"error"`
	}

	errorMsg := fmt.Sprintf("%s: %s", msg, err)
	errorPayload := errorResponse{Error: errorMsg}

	respondWithJSON(w, code, errorPayload)
}
