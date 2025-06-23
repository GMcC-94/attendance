package types

import "github.com/golang-jwt/jwt/v5"

type Credentials struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserID int `json:"userID"`
	jwt.RegisteredClaims
}
