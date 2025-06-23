package main

import (
	"encoding/json"
	"fmt"
	"http-server/internal/auth"
	"net/http"
	"time"
)

type Login struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
type LoginReturn struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
	fmt.Println("there")
	fmt.Println("USER from DB PW:", user.HashedPassword.String)
	fmt.Println("Request PW:", loginReq.Password)
	err = auth.CheckPasswordHash(user.HashedPassword.String, loginReq.Password)
	if err != nil {
		returnJsonError(w, "Incorrect email or password", 401)
		return
	}
	fmt.Println("where")
	rtn := LoginReturn{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	if err := json.NewEncoder(w).Encode(rtn); err != nil {
		returnJsonError(w, "Error endcoding response to json", 401)
		return
	}
}
