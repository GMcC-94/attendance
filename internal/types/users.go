package types

import (
	"errors"
	"time"
)

var ErrUsernameTaken = errors.New("username already taken")

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
