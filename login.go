package main

import (
	"encoding/json"
	"http-server/internal/auth"
	"net/http"
	"time"
)

type Login struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}
type LoginReturn struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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

	rtn := LoginReturn{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwt,
	}

	if err := json.NewEncoder(w).Encode(rtn); err != nil {
		returnJsonError(w, "Error endcoding response to json", 401)
		return
	}
}
