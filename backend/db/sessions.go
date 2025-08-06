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
	User(token string) (*types.User, error)
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

func (p *PostgresSessionsStore) User(token string) (*types.User, error) {
	tokenHash := p.hash(token)

	var user types.User
	row := p.DB.QueryRow(`
	SELECT users.id,
	users.username,
	users.password_hash
	FROM sessions
	JOIN users on users.id = sessions.user_id
	WHERE sessions.token_hash = $1`, tokenHash)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user session: %w", err)
	}
	return &user, nil
}

func (p *PostgresSessionsStore) Delete(token string) error {
	tokenHash := p.hash(token)

	_, err := p.DB.Exec(`
	DELETE FROM sessions
	WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
func (p *PostgresSessionsStore) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
