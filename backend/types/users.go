package types

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username" validate:"required,alphanum,min=3,max=50"`
	PasswordHash string `json:"passwordHash" validate:"required,min=6"`
}

type NewUser struct {
	Username string
	Password string
}

// For API input
type Credentials struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=50"`
	Password string `json:"password" validate:"required"`
}

type Claims struct {
	UserID int `json:"userID"`
	jwt.RegisteredClaims
}
