package config

import (
	"encoding/base64"
	"errors"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr      string
	DatabaseURL   string
	LogMode       string
	JWTKey        []byte
	JWTIssuer     string
	JWTAudience   string
	CookieSecure  bool
	AllowedOrigin string
}

func Load() (*Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8083"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	logMode := os.Getenv("LOG_MODE")
	if logMode == "" {
		logMode = "default"
	}

	// JWT_SECRET — base64 от минимум 32 случайных байт, не обычный пароль.
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
	origin := os.Getenv("ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:8082"
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("ALLOWED_ORIGIN must be an http(s) origin without path")
	}

	return &Config{
		HTTPAddr:    addr,
		DatabaseURL: databaseURL,
		LogMode:     logMode,
		JWTKey:      key, JWTIssuer: issuer, JWTAudience: audience, CookieSecure: secure, AllowedOrigin: origin,
	}, nil
}
