package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

func CreateUserHandler(userStore db.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials types.Credentials
		if !helpers.DecodeJSON(r, w, &credentials) {
			return
		}

		if _, err := helpers.ValidateStruct(credentials); err != nil {
			http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := userStore.CreateUser(credentials.Username, credentials.Password)
		if err != nil {
			if errors.Is(err, types.ErrUsernameTaken) {
				helpers.JSONError(w, http.StatusConflict, "Username already taken", nil)
				return
			}
			helpers.JSONError(w, http.StatusInternalServerError, "User registration failed", nil)
			return
		}

		accessToken, err := helpers.GenerateJWT(userID, 15*time.Minute)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]string{
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
		if !helpers.DecodeJSON(r, w, &credentials) {
			return
		}

		if fields, err := helpers.ValidateStruct(credentials); err != nil {
			helpers.JSONError(w, http.StatusUnauthorized, "Invalid Credentials", fields)
			return
		}

		user, err := userStore.AuthenticateUser(credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		accessToken, err := helpers.GenerateJWT(user.ID, 15*time.Minute)
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
		err = rTokenStore.SaveRefreshToken(user.ID, refreshToken, expiresAt)
		if err != nil {
			http.Error(w, "error saving refresh token", http.StatusInternalServerError)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/api/v1/refresh",
			Expires:  expiresAt,
		})

		helpers.WriteJSON(w, http.StatusOK, map[string]string{
			"message":      "User successfully logged in",
			"access_token": accessToken,
		})
	}
}
