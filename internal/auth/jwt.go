package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn * time.Second)),
		Subject:   userID.String(),
	})
	jwt, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	uuidID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	rawToken := headers.Get("Authorization")
	if len(rawToken) == 0 {
		return "", fmt.Errorf("authorization token not present")
	}
	token := strings.TrimPrefix(rawToken, "Bearer ")
	return token, nil
}

func MakeRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	str := hex.EncodeToString(bytes)
	return str, nil
}
