package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gmcc94/attendance-go/config"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPass,
		config.AppConfig.DBName,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
