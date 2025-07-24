package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
)

func RefreshTokenHandler(rTokenStore *db.PostgresRefreshTokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RefreshToken string `json:"refreshToken"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userID, err := rTokenStore.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
			return
		}

		newAccessToken, err := helpers.GenerateJWT(userID, 15*time.Minute)
		if err != nil {
			http.Error(w, "Failed to generate new token", http.StatusInternalServerError)
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]string{
			"access_token": newAccessToken,
		})
	}
}
