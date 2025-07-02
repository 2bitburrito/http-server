package auth

import (
	"fmt"
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

func TestMakeRefreshToken(t *testing.T) {
	tokens := make([]string, 5)

	for i := range len(tokens) {
		token, err := MakeRefreshToken()
		if err != nil {
			t.Fatalf("failed to make token: %s", err)
		}
		tokens[i] = token
	}

	for i, token := range tokens {
		for j := i + 1; j < len(tokens); j++ {
			if token == tokens[j] {
				fmt.Println(tokens)
				t.Fatalf("Found matching tokens")
			}
		}
	}
}
