package main

import (
	"database/sql"
	"encoding/json"
	"http-server/internal/auth"
	"http-server/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) addUser(w http.ResponseWriter, req *http.Request) {
	var r Request
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		log.Println("Error decoding json in AddUser")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Received email:", r.Email)

	hashedPw, err := auth.HashPassword(r.Password)
	if err != nil {
		returnJsonError(w, "ERROR HASHING", 500)
	}
	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email: r.Email,
		HashedPassword: sql.NullString{
			Valid:  true,
			String: hashedPw,
		},
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	dat, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshalling user in addUser")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(201)
	w.Write(dat)
}
