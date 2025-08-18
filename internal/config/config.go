package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
	Environment   string
	Database      DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	dbConfig := DatabaseConfig{
		Host:     getEnv("DB_HOST", ""),
		Port:     getEnv("DB_PORT", ""),
		User:     getEnv("DB_USER", ""),
		Password: getEnv("DB_PASS", ""),
		Name:     getEnv("DB_NAME", ""),
		SSLMode:  getEnv("DB_SSL", ""),
	}

	cfg := &Config{
		DatabaseURL:   buildDatabaseURL(dbConfig),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		Environment:   getEnv("ENVIRONMENT", "development"),
		Database:      dbConfig,
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildDatabaseURL(db DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

func (c *Config) validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASS is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.Database.SSLMode == "" {
		return fmt.Errorf("DB_SSL is required")
	}
	return nil
}
