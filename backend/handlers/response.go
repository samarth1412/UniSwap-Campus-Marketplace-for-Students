package handlers

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, apiResponse{
		Success: true,
		Data:    data,
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, apiResponse{
		Success: false,
		Error:   message,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
