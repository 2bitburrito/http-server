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

func (cfg *apiConfig) updateUser(w http.ResponseWriter, req *http.Request) {
	var r Request
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		returnJsonError(w, "Error decoding json in update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnJsonError(w, "Error Getting Token", 401)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		returnJsonError(w, "Couldnt' validate JWT", 401)
	}
	pw, err := auth.HashPassword(r.Password)
	if err != nil {
		returnJsonError(w, "Error hashing password", 500)
		return
	}
	user, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		ID: userID,
		HashedPassword: sql.NullString{
			Valid:  true,
			String: pw,
		},
		Email: r.Email,
	})
	if err != nil {
		returnJsonError(w, "Error setting new email/password in db: "+err.Error(), 500)
		return
	}
	rtnUser := User{
		ID:    user.ID,
		Email: user.Email,
	}
	if err := json.NewEncoder(w).Encode(rtnUser); err != nil {
		returnJsonError(w, "Couldn't set json response: "+err.Error(), 500)
	}
}
