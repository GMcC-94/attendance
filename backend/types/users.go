package types

import "github.com/golang-jwt/jwt/v5"

type Credentials struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Claims struct {
	UserID int `json:"userID"`
	jwt.RegisteredClaims
}
