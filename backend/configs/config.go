package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	OpenAIKey   string
	S3Bucket    string
	S3Region    string
	S3Endpoint  string
	AWSKey      string
	AWSSecret   string
	Environment string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	return &Config{
		Port:        getEnv("PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/resume_screener?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
		S3Bucket:    getEnv("S3_BUCKET", "resume-screener"),
		S3Region:    getEnv("S3_REGION", "auto"),
		S3Endpoint:  getEnv("S3_ENDPOINT", ""),
		AWSKey:      getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecret:   getEnv("AWS_SECRET_ACCESS_KEY", ""),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
