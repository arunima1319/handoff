package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {

	uuidTest := uuid.New()

	expiration, err := time.ParseDuration("300s")
	if err != nil {
		t.Errorf("Could not get expiration time: %s", err)
	}

	tokenString, err := MakeJWT(uuidTest, "my-random-secret-lala", expiration)
	if err != nil {
		t.Errorf("Could not make JWT: %s", err)
	}

	actualUUID, err := ValidateJWT(tokenString, "my-random-secret-lala")
	if err != nil {
		t.Errorf("Could not validate JWT: %s", err)
	}

	if uuidTest != actualUUID {
		t.Errorf("UUIDs do not match")
	}

}
