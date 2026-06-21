package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Hold config variables
type Config struct {
	DatabaseURL   string
	RedisURL      string
	JwtKey        string
	MailgunAPIKey string
	MailgunDomain string
}

// Load environment variables
func Load() (*Config, error) {
	// Try to inject from .env file if available
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	dbURL, found := os.LookupEnv("POSTGRES_DB_URL")
	if !found {
		return nil, errors.New("Missing POSTGRES_DB_URL environment variable")
	}
	redisURL, found := os.LookupEnv("REDIS_URL")
	if !found {
		return nil, errors.New("Missing REDIS_URL environment variable")
	}
	jwtKey, found := os.LookupEnv("JWT_SECRET")
	if !found {
		return nil, errors.New("Missing JWT_SECRET environment variable")
	}
	mailgunApiKey, found := os.LookupEnv("MAILGUN_API_KEY")
	if !found {
		return nil, errors.New("Missing MAILGUN_API_KEY environment variable")
	}
	mailgunDomain, found := os.LookupEnv("MAILGUN_DOMAIN")
	if !found {
		return nil, errors.New("Missing MAILGUN_DOMAIN environment variable")
	}

	return &Config{
		DatabaseURL:   dbURL,
		RedisURL:      redisURL,
		JwtKey:        jwtKey,
		MailgunAPIKey: mailgunApiKey,
		MailgunDomain: mailgunDomain,
	}, nil
}
