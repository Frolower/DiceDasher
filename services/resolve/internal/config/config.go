package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	LogMode     string
}

func Load() (*Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	logMode := os.Getenv("LOG_MODE")
	if logMode == "" {
		logMode = "default"
	}

	return &Config{
		HTTPAddr:    addr,
		DatabaseURL: databaseURL,
		LogMode:     logMode,
	}, nil
}
