package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func reqJSONError(err error) (code int, msg string) {
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return http.StatusBadRequest, "Malformed JSON: "
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return http.StatusBadRequest, "Incorrect fields: "
	}

	return http.StatusBadRequest, "Invalid request body: "

}
