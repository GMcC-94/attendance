package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
}

var AppConfig *Config

func LoadConfig() {
	_ = godotenv.Load()

	AppConfig = &Config{
		JWTSecret: getEnv("JWT_SECRET"),
		DBHost:    getEnv("DB_HOST"),
		DBPort:    getEnv("DB_PORT"),
		DBUser:    getEnv("DB_USER"),
		DBPass:    getEnv("DB_PASS"),
		DBName:    getEnv("DB_NAME"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env variable %s not set", key)
	}

	return val
}
