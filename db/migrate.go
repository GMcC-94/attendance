package db

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

func RunMigrationsUp(db *sql.DB) {
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}

func RunMigrationsDown(db *sql.DB) {
	if err := goose.Down(db, "migrations"); err != nil {
		log.Fatalf("Failed to rollback migration: %v", err)
	}
	log.Println("Migration rolled back successfully")
}

func RunMigrationStatus(db *sql.DB) {
	if err := goose.Status(db, "migrations"); err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}

}
