package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "handoff-access",
			Subject:   userID.String(),
			IssuedAt:  &jwt.NumericDate{time.Now().UTC()},
			ExpiresAt: &jwt.NumericDate{time.Now().Add(expiresIn).UTC()},
		})

	accessToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("create token string: %w", err)
	}
	return accessToken, nil
}

func ValidateJWT(accessToken, tokenSecret string) (uuid.UUID, error) {

	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(t *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse jwt: %w", err)
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("get subject from token: %w", err)
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {

		return uuid.Nil, fmt.Errorf("invalid user id string: %w", err)
	}

	return userID, nil
}
