package main

import (
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
	defer req.Body.Close()
	err := decoder.Decode(&r)
	if err != nil {
		returnJsonError(w, "Error Decoding from json"+err.Error(), http.StatusInternalServerError)
		return
	}
	length := len(r.Body)
	if length > 140 {
		returnJsonError(w, "Chirp is too long", 400)
		return
	}

	cleanBody := cleanChirp(r.Body)
	dbChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:      cleanBody,
		UserID:    r.UserID,
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
		returnJsonError(w, "error getting all chrips from db"+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(allChirps)
	w.WriteHeader(200)
}

func (cfg *apiConfig) getSingleChirp(w http.ResponseWriter, req *http.Request) {
	var chirpID uuid.UUID
	chirpID = req.PathValue("chirp-id")
	chirpUUID := uuid
	if !ok {
		returnJsonError(w, "Error type casting to uuid", http.StatusInternalServerError)
	}

	cfg.dbQueries.GetSingleChirp(req.Context(), chirpID)
}
