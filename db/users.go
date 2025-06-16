package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SignupUser(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error: Hashing Password %v", hashedPassword)
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hashedPassword))
	return err
}

func AuthenticateUser(db *sql.DB, username, password string) (bool, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}
