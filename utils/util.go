package utils

import (
	"encoding/json"
	"net/http"
)

// Message -> Function to build json message
func Message(success bool, message string) map[string]interface{} {
	return map[string]interface{}{"success": success, "message": message}
}

// Respond -> Function to wrap request
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// RespondWithStatusBadRequest -> Function to respond with a status 400
func RespondWithStatusBadRequest(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(data)
}

// RespondWithStatusInternalServerError -> Function to respond with a status 500
func RespondWithStatusInternalServerError(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(data)
}

// RespondWithStatusUnauthorized -> Function to respond with a status unauthorized
func RespondWithStatusUnauthorized(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(data)
}

// RespondWithStatusNotFound -> Function to respond with when a record is not found
func RespondWithStatusNotFound(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(data)
}
