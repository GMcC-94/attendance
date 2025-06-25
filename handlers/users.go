package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

func SignupHandler(userStore db.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials types.Credentials
		json.NewDecoder(r.Body).Decode(&credentials)

		if err := helpers.ValidateStruct(credentials); err != nil {
			http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := userStore.SignupUser(credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "user registration failed", http.StatusInternalServerError)
			return
		}

		accessToken, err := helpers.GenerateJWT(userID, 15*time.Minute)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":      "User signed up successfully",
			"access_token": accessToken,
		})
	}
}

func LoginHandler(
	userStore db.UserStore,
	rTokenStore db.RefreshTokenStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials types.Credentials
		json.NewDecoder(r.Body).Decode(&credentials)

		if err := helpers.ValidateStruct(credentials); err != nil {
			http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := userStore.AuthenticateUser(credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		accessToken, err := helpers.GenerateJWT(userID, 15*time.Minute)
		if err != nil {
			http.Error(w, "error generating token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := helpers.GenerateRandomToken()
		if err != nil {
			http.Error(w, "error generating refresh token", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		err = rTokenStore.SaveRefreshToken(userID, refreshToken, expiresAt)
		if err != nil {
			http.Error(w, "error saving refresh token", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken})
	}
}
