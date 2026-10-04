package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {

	authString, ok := headers["Authorization"]
	if !ok {
		return "", errors.New("No authorization header found")
	}
	credentials := strings.Split(authString[0], " ")
	if len(credentials) < 2 {
		return "", errors.New("Access token missing in value")
	}
	return credentials[1], nil
}
