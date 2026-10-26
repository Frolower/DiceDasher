package config

import (
	"diceDasher/pkg/httputil"
	"errors"
	"os"
)

type Config struct {
	HTTPAddr       string
	DatabaseURL    string
	LogMode        string
	AllowedOrigins []string
}

func Load() (*Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8082"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	logMode := os.Getenv("LOG_MODE")
	if logMode == "" {
		logMode = "default"
	}

	origins, err := httputil.ParseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"), "http://localhost:8888,http://localhost:8000")
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTPAddr:       addr,
		DatabaseURL:    databaseURL,
		LogMode:        logMode,
		AllowedOrigins: origins,
	}, nil
}
