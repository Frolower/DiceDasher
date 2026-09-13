package auth

import (
	"time"

	"github.com/google/uuid"
)

// CreateInput user creation input
type CreateInput struct {
	Username string
	Email    string
	Password string
}

// CreateRecord user creation db input
type CreateRecord struct {
	Username     string
	Email        string
	PasswordHash string
}

// Created user creation response
type Created struct {
	ID uuid.UUID
}

type LoginInput struct {
	Username  string
	Password  string
	UserAgent string
}

// Tokens login/refresh response
type Tokens struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

// Credentials login action verification struct
type Credentials struct {
	UserID       uuid.UUID
	PasswordHash string
}

// Session описывает долгоживущую сессию. Access JWT в БД не хранится:
// он содержит sid и собственный exp; одна сессия может выдавать много JWT.
type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	CreatedAt        time.Time
	ExpiresAt        time.Time
	// nil означает SQL NULL: сессия ещё не отозвана.
	RevokedAt *time.Time
	UserAgent string
}
