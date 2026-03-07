package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (*Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	return &Config{
		HTTPAddr:    addr,
		DatabaseURL: databaseURL,
	}, nil
}
