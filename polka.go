package main

import (
	"encoding/json"
	"http-server/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

type Webhook struct {
	Event string      `json:"event"`
	Data  WebhookData `json:"data"`
}
type WebhookData struct {
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, req *http.Request) {
	receivedKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		returnJsonError(w, "Error getting api key: "+err.Error(), 401)
		return
	}
	if receivedKey != cfg.polkaKey {
		w.WriteHeader(401)
		return
	}

	var data Webhook
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		returnJsonError(w, "Couldn't decode json in handlePolkaWebhook", 500)
		return
	}
	defer req.Body.Close()

	if data.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	if err := cfg.dbQueries.UpgradeUserToRed(req.Context(), data.Data.UserID); err != nil {
		returnJsonError(w, "Couldn't change user red status in db: "+err.Error(), 404)
		return
	}
	w.WriteHeader(204)
}
