package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gmcc94/attendance-go/config"
)

type APIResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Error   string            `json:"error,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fallback := map[string]string{"error": "failed to encode response"}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(fallback)
	}
}

func JSONError(w http.ResponseWriter, status int, message string, fields map[string]string) {
	resp := APIResponse{
		Success: false,
		Error:   message,
	}

	if config.IsDev && len(fields) > 0 {
		resp.Fields = fields
	}
	WriteJSON(w, status, resp)
}

func DecodeJSON(r *http.Request, w http.ResponseWriter, v interface{}) bool {
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		log.Printf("JSON decode error on %s: %v | Body: %s", r.URL.Path, err, string(bodyBytes))
		JSONError(w, http.StatusBadRequest, "Invalid input", nil)
		return false
	}
	return true
}

func JSONSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	WriteJSON(w, status, resp)
}
