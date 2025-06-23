package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateAndValidateJWT(t *testing.T) {
	tokenSecret := "laskdjflsjkdf"
	userID := uuid.New()
	expiresIn := time.Minute
	jwt, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Error creating JWT: %s", err.Error())
	}
	returnedID, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.Errorf("Error validating JWT: %s", err.Error())
	}
	if userID != returnedID {
		t.Errorf("User ID's do not match after JWT creation")
	}
}
