package main

const msgServerError = "internal server error"
const msgFailedAuthentication = "authentication failed"
const msgLoginError = "invalid user email or password"
const msgInvalidRequestBody = "invalid request body"
const msgForbiddenError = "user forbidden from carrying out this action"

/*
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
*/
