package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type request struct {
	Body string `json:"body"`
}

type validResp struct {
	CleanedBody string `json:"cleaned_body"`
}
type errorResp struct {
	Error string `json:"error"`
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
	var r request
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(&r)
	if err != nil {
		respBody := errorResp{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s\n", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
	length := len(r.Body)
	if length > 140 {
		respBody := errorResp{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s\n", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	cleanBody := cleanChirp(r.Body)

	respBody := validResp{
		CleanedBody: cleanBody,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
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
