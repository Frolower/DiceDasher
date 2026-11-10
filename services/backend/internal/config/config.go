package config

import (
	"diceDasher/pkg/httputil"
	"encoding/base64"
	"encoding/json/v2"
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
	HealthTargets  map[string]string
	RoomGRPCAddr   string
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

	targets := map[string]string{
		"room":      "grpc://localhost:50051",
		"character": "http://localhost:8082/health",
		"resolve":   "http://localhost:8081/health",
	}
	if value := os.Getenv("HEALTH_TARGETS"); value != "" {
		targets = nil
		if err := json.Unmarshal([]byte(value), &targets); err != nil || targets == nil {
			return nil, errors.New("HEALTH_TARGETS must be a JSON object mapping service names to URLs")
		}
	}

	roomAddr := os.Getenv("ROOM_GRPC_ADDR")
	if roomAddr == "" {
		roomAddr = "localhost:50051"
	}
	return &Config{
		RoomGRPCAddr:  roomAddr,
		HTTPAddr:      addr,
		DatabaseURL:   databaseURL,
		LogMode:       logMode,
		HealthTargets: targets,
		JWTKey:        key, JWTIssuer: issuer, JWTAudience: audience, CookieSecure: secure, AllowedOrigins: origins,
	}, nil
}
