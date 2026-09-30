package auth

import (
	"log"
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
		log.Printf("Error in creating token string: %v", err)
		return "", err
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
		log.Printf("Could not parse token: %v", err)
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Could not get subject from token")
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("Could not convert user ID string to UUID: %v", err)
		return uuid.Nil, err
	}

	return userID, nil
}
