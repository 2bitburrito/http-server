package main

import (
	"database/sql"
	"encoding/json"
	"http-server/internal/database"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type request struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type validResp struct {
	CleanedBody string `json:"cleaned_body"`
}
type errorResp struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) postChirp(w http.ResponseWriter, req *http.Request) {
	var r request
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&r)
	if err != nil {
		returnJsonError(w, "Error Decoding from json"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	length := len(r.Body)
	if length > 140 {
		returnJsonError(w, "Chirp is too long", 400)
		return
	}

	cleanBody := cleanChirp(r.Body)
	userID := req.Context().Value("UserID")
	uuidID, _ := uuid.Parse(userID.(string))
	dbChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:      cleanBody,
		UserID:    uuidID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		returnJsonError(w, "Error inserting chrip into database"+err.Error(), http.StatusInternalServerError)
	}

	dat, err := json.Marshal(dbChirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(dat)
}

func cleanChirp(body string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	words := strings.Fields(body)
	cleanStringsSlice := make([]string, len(words))

	for i, word := range words {
		replaced := false
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				replaced = true
				cleanStringsSlice[i] = "****"
			}
		}
		if !replaced {
			cleanStringsSlice[i] = word
		}
	}
	cleanStrings := strings.Join(cleanStringsSlice, " ")
	return cleanStrings
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, req *http.Request) {
	allChirps, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		returnJsonError(w, "error getting all chrips from db: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(allChirps)
	w.WriteHeader(200)
}

func (cfg *apiConfig) getSingleChirp(w http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("id")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		returnJsonError(w, "Error type casting to uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	chirp, err := cfg.dbQueries.GetSingleChirp(req.Context(), chirpUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			returnJsonError(w, "No Chirp Found: "+err.Error(), http.StatusNotFound)
			return
		}
		returnJsonError(w, "Error getting chirp from db: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(chirp); err != nil {
		returnJsonError(w, "Error encoding chirp to json: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
