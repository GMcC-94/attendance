package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

func CreateUserHandler(userStore db.UserStore, sessionStore db.SessionStore) http.HandlerFunc {
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

		session, err := sessionStore.CreateSession(userID)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/login", http.StatusNotFound)
			return
		}
		setCookie(w, CookieSession, session.Token)
		http.Redirect(w, r, "/attendance", http.StatusFound)

		helpers.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "User successfully created",
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
