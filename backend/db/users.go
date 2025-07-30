package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gmcc94/attendance-go/types"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	CreateUser(username, password string) (int, error)
	AuthenticateUser(username, password string) (*types.User, error)
}

type PostgresUserStore struct {
	DB *sql.DB
}

func (p *PostgresUserStore) CreateUser(username, password string) (int, error) {
	username = strings.ToLower(username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error: Hashing Password %v", hashedPassword)
	}

	var userID int
	err = p.DB.QueryRow("INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", username, string(hashedPassword)).
		Scan(&userID)
	if err != nil {
		var pgError interface {
			SQLState() string
		}
		if errors.As(err, &pgError) {
			if pgError.SQLState() == pgerrcode.UniqueViolation {
				return 0, types.ErrUsernameTaken
			}
		}
		log.Printf("error: Inserting user %v", err)
		return 0, fmt.Errorf("create user: %w", err)
	}
	return userID, nil
}

func (p *PostgresUserStore) AuthenticateUser(username, password string) (*types.User, error) {
	username = strings.ToLower(username)
	user := types.User{
		Username: username,
	}

	err := p.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = $1", username).Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return &user, nil
}
