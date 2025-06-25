package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	SignupUser(username, password string) (int, error)
	AuthenticateUser(username, password string) (int, error)
}

type PostgresUserStore struct {
	DB *sql.DB
}

func (p *PostgresUserStore) SignupUser(username, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error: Hashing Password %v", hashedPassword)
	}

	var userID int
	err = p.DB.QueryRow("INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", username, string(hashedPassword)).
		Scan(&userID)
	if err != nil {
		log.Printf("error: Inserting user %v", err)
	}
	return userID, nil
}

func (p *PostgresUserStore) AuthenticateUser(username, password string) (int, error) {
	var userID int
	var hashedPassword string
	err := p.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = $1", username).Scan(&userID, &hashedPassword)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return 0, nil
	}

	return userID, nil
}
