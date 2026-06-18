package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	JwtSecret       string
	AuthServiceURL  string
	UserServiceURL  string
	PostsServiceURL string
	FeedServiceURL  string
}

func Load() (*Config, error) {
	// Try to load from .env
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	// Read env values
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return nil, fmt.Errorf("missing PORT env variable")
	}
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return nil, fmt.Errorf("missing JWT_SECRET env variable")
	}
	authServiceURL, ok := os.LookupEnv("AUTH_SERVICE_URL")
	if !ok {
		return nil, fmt.Errorf("missing AUTH_SERVICE_URL env variable")
	}
	userServiceURL, ok := os.LookupEnv("USER_SERVICE_URL")
	if !ok {
		return nil, fmt.Errorf("missing USER_SERVICE_URL env variable")
	}
	postsServiceURL, ok := os.LookupEnv("POSTS_SERVICE_URL")
	if !ok {
		return nil, fmt.Errorf("missing POSTS_SERVICE_URL env variable")
	}
	feedServiceURL, ok := os.LookupEnv("FEED_SERVICE_URL")
	if !ok {
		return nil, fmt.Errorf("missing FEED_SERVICE_URL env variable")
	}

	return &Config{
		Port:            port,
		JwtSecret:       jwtSecret,
		AuthServiceURL:  authServiceURL,
		UserServiceURL:  userServiceURL,
		PostsServiceURL: postsServiceURL,
		FeedServiceURL:  feedServiceURL,
	}, nil
}
