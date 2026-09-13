package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const dummyPasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

func (s *Service) Login(ctx context.Context, input LoginInput) (Tokens, error) {
	if s.users == nil || s.sessions == nil || s.passwords == nil || s.tokens == nil {
		return Tokens{}, ErrAuthNotConfigured
	}
	if err := validateLoginInput(input); err != nil {
		return Tokens{}, InvalidInputError{Err: err}
	}
	credentials, err := s.users.FindUserByUsername(ctx, input.Username)
	if errors.Is(err, ErrUserNotFound) {
		_ = s.passwords.Verify(dummyPasswordHash, input.Password)
		return Tokens{}, ErrUnableToLogin
	}
	if err != nil {
		return Tokens{}, fmt.Errorf("find login user: %w", err)
	}
	if err := s.passwords.Verify(credentials.PasswordHash, input.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return Tokens{}, ErrUnableToLogin
		}
		return Tokens{}, fmt.Errorf("verify password: %w", err)
	}
	if credentials.UserID == uuid.Nil {
		return Tokens{}, errors.New("stored user has no id")
	}
	return s.createSession(ctx, credentials.UserID, input.UserAgent)
}
