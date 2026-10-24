package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserStore interface {
	CreateUser(ctx context.Context, rec CreateRecord) (uuid.UUID, error)
	FindUserByUsername(ctx context.Context, username string) (Credentials, error)
}

type SessionStore interface {
	CreateSession(ctx context.Context, rec Session) error
	FindSessionByRefreshHash(ctx context.Context, hash string) (Session, error)
	FindSessionByID(ctx context.Context, id uuid.UUID) (Session, error)

	RotateRefreshToken(
		ctx context.Context,
		sessionID uuid.UUID,
		oldHash string,
		newHash string,
		now time.Time,
	) error

	RevokeSessionByRefreshHash(ctx context.Context, hash string, now time.Time) error
}

type Service struct {
	users     UserStore
	sessions  SessionStore
	passwords PasswordManager
	tokens    TokenManager
	now       func() time.Time
}

// NewService creates a service for the existing registration flow.
// For login, main supplies sessions and JWT support through NewServiceWithDependencies.
func NewService(users UserStore) *Service {
	return NewServiceWithDependencies(users, nil, nil, nil)
}

// NewServiceWithPasswords allows replacing bcrypt in registration tests.
func NewServiceWithPasswords(users UserStore, passwords PasswordManager) *Service {
	return NewServiceWithDependencies(users, nil, passwords, nil)
}

// All dependencies are stored once when main wires up the application.
// A nil passwords dependency selects the default bcrypt implementation; sessions and tokens are optional for Register.
func NewServiceWithDependencies(users UserStore, sessions SessionStore, passwords PasswordManager, tokens TokenManager) *Service {
	if passwords == nil {
		passwords = NewBcryptHasher()
	}
	return &Service{users: users, sessions: sessions, passwords: passwords, tokens: tokens, now: time.Now}
}
