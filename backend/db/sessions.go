package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/gmcc94/attendance-go/rand"
	"github.com/gmcc94/attendance-go/types"
)

type SessionStore interface {
	CreateSession(userID int) (*types.Session, error)
}

type PostgresSessionsStore struct {
	DB            *sql.DB
	BytesPerToken int
}

func (p *PostgresSessionsStore) CreateSession(userID int) (*types.Session, error) {
	bytesPerToken := p.BytesPerToken
	if bytesPerToken < types.MinBytesPerToken {
		bytesPerToken = types.MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	session := types.Session{
		UserID:    userID,
		Token:     token,
		TokenHash: p.hash(token),
	}
	row := p.DB.QueryRow(`
	INSERT INTO sessions (user_id, token_hash)
	VALUES ($1, $2) ON CONFLICT (user_id) DO
	UPDATE
	SET token_hash = $2
	RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &session, nil
}

func (p *PostgresSessionsStore) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
