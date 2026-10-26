package config

import (
	"diceDasher/pkg/httputil"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr       string
	DatabaseURL    string
	LogMode        string
	JWTKey         []byte
	JWTIssuer      string
	JWTAudience    string
	CookieSecure   bool
	AllowedOrigins []string
}

func Load() (*Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8000"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	logMode := os.Getenv("LOG_MODE")
	if logMode == "" {
		logMode = "default"
	}

	// JWT_SECRET is the base64 encoding of at least 32 random bytes, not a regular password.
	key, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_SECRET"))
	if err != nil || len(key) < 32 {
		return nil, errors.New("JWT_SECRET must be base64 encoding of at least 32 random bytes")
	}
	secure := true
	if value := os.Getenv("COOKIE_SECURE"); value != "" {
		secure, err = strconv.ParseBool(value)
		if err != nil {
			return nil, errors.New("COOKIE_SECURE must be true or false")
		}
	}
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "dicedasher-backend"
	}
	audience := os.Getenv("JWT_AUDIENCE")
	if audience == "" {
		audience = "dicedasher-api"
	}
	origins, err := httputil.ParseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"), "http://localhost:8888,http://localhost:5174")
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTPAddr:    addr,
		DatabaseURL: databaseURL,
		LogMode:     logMode,
		JWTKey:      key, JWTIssuer: issuer, JWTAudience: audience, CookieSecure: secure, AllowedOrigins: origins,
	}, nil
}
