package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func returnJsonError(w http.ResponseWriter, msg string, code int) {
	respBody := errorResp{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err)
		return
	}
	log.Println(msg)
	w.WriteHeader(code)
	w.Write(dat)
}
