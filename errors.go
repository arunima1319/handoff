package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const msgServerError = "internal server error"

func reqJSONError(err error) (code int, msg string) {
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		log.Printf("Malformed JSON: %s", err)
		return http.StatusBadRequest, "JSON is malformed"
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		log.Printf("Incorrect field types: %s", err)
		return http.StatusBadRequest, "Field types are incorrect"
	}

	return http.StatusBadRequest, "Invalid request body"

}
