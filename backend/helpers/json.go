package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fallback := map[string]string{"error": "failed to encode response"}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(fallback)
	}
}

func JSONError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func DecodeJSON(r *http.Request, w http.ResponseWriter, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		log.Printf("JSON decode error: %v", err)
		JSONError(w, http.StatusBadRequest, "Invalid input")
		return false
	}
	return true
}
