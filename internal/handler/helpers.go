package handler

import (
	"encoding/json"
	"net/http"
)

// writeJSON sends a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Encoding error after WriteHeader — nothing we can do gracefully.
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
	}
}

// writeError sends a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	//nolint:errcheck
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
