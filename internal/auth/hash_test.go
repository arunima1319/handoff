package auth

import "testing"

func TestHashPassword(t *testing.T) {

	password := "randomPassword129@()"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Errorf("could not hash password: %v", err)
	}

	match, _ := CheckPasswordHash(hashed, password)
	if match != true {
		t.Errorf("password does not match")
	}

}
