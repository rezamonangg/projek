package common

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	WriteJSON(w, status, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, APIResponse{
		Success: false,
		Error:   &APIError{Code: code, Message: message},
	})
}
