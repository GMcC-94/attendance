package db

import (
	"database/sql"
	"time"
)

type RefreshTokenStore interface {
	SaveRefreshToken(userID int, token string, expiresAt time.Time) error
	ValidateRefreshToken(token string) (int, error)
}

type PostgresRefreshTokenStore struct {
	DB *sql.DB
}

func (p *PostgresRefreshTokenStore) SaveRefreshToken(userID int, token string, expiresAt time.Time) error {
	_, err := p.DB.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, token, expiresAt)
	return err
}

func (p *PostgresRefreshTokenStore) ValidateRefreshToken(token string) (int, error) {
	var userID int
	var expiresAt time.Time

	err := p.DB.QueryRow("SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1", token).
		Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, sql.ErrNoRows
	}

	return userID, nil
}
