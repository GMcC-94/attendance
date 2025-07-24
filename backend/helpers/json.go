package helpers

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error: failed to encode response", http.StatusInternalServerError)
	}
}

func JSONError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error:": message})
}

func DecodeJSON(r *http.Request, w http.ResponseWriter, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		JSONError(w, http.StatusBadRequest, "Invlaid input")
		return false
	}
	return true
}
