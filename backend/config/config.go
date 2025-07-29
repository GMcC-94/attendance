package config

import (
	"log"
	"os"

	"github.com/gmcc94/attendance-go/types"
	"github.com/joho/godotenv"
)

var AppConfig *types.Config
var IsDev = os.Getenv("APP_ENV") == "development"

func LoadConfig() {
	_ = godotenv.Load()

	AppConfig = &types.Config{
		JWTSecret:          getEnv("JWT_SECRET"),
		DBHost:             getEnv("DB_HOST"),
		DBPort:             getEnv("DB_PORT"),
		DBUser:             getEnv("DB_USER"),
		DBPass:             getEnv("DB_PASS"),
		DBName:             getEnv("DB_NAME"),
		AWSRegion:          getEnv("AWS_REGION"),
		AWSBucketName:      getEnv("AWS_BUCKET_NAME"),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY"),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env variable %s not set", key)
	}

	return val
}
