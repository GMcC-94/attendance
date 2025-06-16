package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	userDB "github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

func SignupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials types.Credentials
		json.NewDecoder(r.Body).Decode(&credentials)

		if err := helpers.ValidateStruct(credentials); err != nil {
			http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
			return
		}

		err := userDB.SignupUser(db, credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "user registration failed", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("User signed up successfully"))
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials types.Credentials
		json.NewDecoder(r.Body).Decode(&credentials)

		if err := helpers.ValidateStruct(credentials); err != nil {
			http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
			return
		}

		valid, err := userDB.AuthenticateUser(db, credentials.Username, credentials.Password)
		if err != nil || !valid {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := helpers.GenerateJWT(credentials.Username)
		if err != nil {
			http.Error(w, "error generating token", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Login Successful"))
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
