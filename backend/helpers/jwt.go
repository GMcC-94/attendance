package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/gmcc94/attendance-go/config"
	"github.com/gmcc94/attendance-go/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int, expiry time.Duration) (string, error) {
	claims := types.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	return tokenString, err
}

func ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})
	if err != nil {
		return 0, nil
	}

	claims, ok := token.Claims.(*types.Claims)
	if !ok || token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}

func GenerateRandomToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
