package main

import (
	"encoding/json"
	"fmt"
	"http-server/internal/auth"
	"http-server/internal/database"
	"net/http"
	"time"
)

type Login struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}
type LoginReturn struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}
type RefreshTokenReturn struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) login(w http.ResponseWriter, req *http.Request) {
	var loginReq Login
	err := json.NewDecoder(req.Body).Decode(&loginReq)
	if err != nil {
		returnJsonError(w, "Error decoding json"+err.Error(), 401)
		return
	}
	user, err := cfg.dbQueries.GetUserByEmaill(req.Context(), loginReq.Email)
	if err != nil {
		returnJsonError(w, "error getting user: "+err.Error(), 401)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword.String, loginReq.Password)
	if err != nil {
		returnJsonError(w, "Incorrect email or password", 401)
		return
	}
	jwtExpiry := loginReq.ExpiresInSeconds
	if jwtExpiry == 0 || jwtExpiry > 3600 {
		jwtExpiry = 3600
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(jwtExpiry))
	if err != nil {
		returnJsonError(w, "Error creating jwt", 500)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		returnJsonError(w, "Error creating jwt", 500)
		return
	}
	// 60 day expiry
	refreshTokenExpiry := time.Now().Add(60 * 24 * time.Hour)
	err = cfg.dbQueries.AddRefreshToken(req.Context(), database.AddRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: refreshTokenExpiry,
	})
	if err != nil {
		returnJsonError(w, "Error Adding Refresh Token to db: "+err.Error(), 401)
		return
	}

	rtn := LoginReturn{
		Id:           user.ID.String(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        jwt,
		RefreshToken: refreshToken,
	}

	if err := json.NewEncoder(w).Encode(rtn); err != nil {
		returnJsonError(w, "Error endcoding response to json", 401)
		return
	}
}

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnJsonError(w, "error while getting token from Headers: "+err.Error(), 401)
		return
	}
	fmt.Printf("received token %s\n", refreshToken)

	tokenRow, err := cfg.dbQueries.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		returnJsonError(w, "error while retreiving token row: "+err.Error(), 401)
		return
	}

	fmt.Printf("stored token %+v\n", tokenRow)

	if tokenRow.ExpiresAt.Before(time.Now()) ||
		tokenRow.RevokedAt.Valid {
		w.WriteHeader(401)
		return
	}

	jwt, err := auth.MakeJWT(tokenRow.UserID, cfg.tokenSecret, time.Duration(3600))
	if err != nil {
		returnJsonError(w, "Error creating jwt", 500)
		return
	}
	returnData := RefreshTokenReturn{
		Token: jwt,
	}
	if err := json.NewEncoder(w).Encode(returnData); err != nil {
		returnJsonError(w, "Error encoding jwt to json: "+err.Error(), 500)
		return
	}
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnJsonError(w, "error while getting token from Headers: "+err.Error(), 401)
		return
	}
	if err := cfg.dbQueries.RevokeToken(req.Context(), refreshToken); err != nil {
		returnJsonError(w, "error while revoking token: "+err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}
