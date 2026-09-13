package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
	"time"
)

type loginUsers struct {
	UserStore
	credentials Credentials
	err         error
}

func (u loginUsers) FindUserByUsername(context.Context, string) (Credentials, error) {
	return u.credentials, u.err
}

type sessionCapture struct {
	SessionStore
	records []Session
	err     error
}

func (s *sessionCapture) CreateSession(_ context.Context, rec Session) error {
	s.records = append(s.records, rec)
	return s.err
}

type brokenTokens struct {
	TokenManager
	generateErr, issueErr error
}

func (t brokenTokens) GenerateRefresh() (string, error) {
	if t.generateErr != nil {
		return "", t.generateErr
	}
	return t.TokenManager.GenerateRefresh()
}
func (t brokenTokens) IssueAccess(u, s uuid.UUID, exp time.Time) (string, error) {
	if t.issueErr != nil {
		return "", t.issueErr
	}
	return t.TokenManager.IssueAccess(u, s, exp)
}

func TestLoginCreatesSession(t *testing.T) {
	passwords := BcryptHasher{Cost: bcrypt.MinCost}
	// Существующий пароль может не удовлетворять сегодняшним правилам регистрации.
	hash, _ := passwords.Hash("old")
	id := uuid.New()
	store := &sessionCapture{}
	tokens, _ := NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
	service := NewServiceWithDependencies(loginUsers{credentials: Credentials{UserID: id, PasswordHash: hash}}, store, passwords, tokens)
	result, err := service.Login(context.Background(), LoginInput{Username: "Alice", Password: "old", UserAgent: "test agent"})
	if err != nil {
		t.Fatal(err)
	}
	if len(store.records) != 1 {
		t.Fatal("session was not saved exactly once")
	}
	rec := store.records[0]
	claims, err := tokens.VerifyAccess(result.AccessToken)
	if err != nil || claims.UserID != id || claims.SessionID != rec.ID {
		t.Fatalf("wrong claims: %v %v", claims, err)
	}
	if rec.RefreshTokenHash != tokens.HashRefresh(result.RefreshToken) || rec.RefreshTokenHash == result.RefreshToken || rec.RevokedAt != nil || rec.UserAgent != "test agent" {
		t.Fatal("wrong session data")
	}
	if rec.ExpiresAt.Sub(rec.CreatedAt) != SessionLifetime || result.AccessExpiresAt.Sub(rec.CreatedAt) != AccessLifetime || !result.RefreshExpiresAt.Equal(rec.ExpiresAt) {
		t.Fatal("wrong lifetimes")
	}
	second, err := service.Login(context.Background(), LoginInput{Username: "Alice", Password: "old"})
	if err != nil || second.RefreshToken == result.RefreshToken || store.records[0].ID == store.records[1].ID {
		t.Fatal("login reused a session/token")
	}
}

func TestLoginFailuresDoNotReturnTokens(t *testing.T) {
	passwords := BcryptHasher{Cost: bcrypt.MinCost}
	hash, _ := passwords.Hash("StrongPass1")
	manager, _ := NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
	failure := errors.New("internal failure")
	for _, tc := range []struct {
		name                                      string
		userErr                                   error
		password, hash                            string
		saveErr, generateErr, issueErr, errorWant error
		saves                                     int
	}{
		{name: "wrong password", password: "WrongPass1", hash: hash, errorWant: ErrUnableToLogin},
		{name: "unknown user", password: "StrongPass1", userErr: ErrUserNotFound, errorWant: ErrUnableToLogin},
		{name: "read failure", password: "StrongPass1", userErr: failure, errorWant: failure},
		{name: "save failure", password: "StrongPass1", hash: hash, saveErr: failure, errorWant: failure, saves: 1},
		{name: "random failure", password: "StrongPass1", hash: hash, generateErr: failure, errorWant: failure},
		{name: "sign failure", password: "StrongPass1", hash: hash, issueErr: failure, errorWant: failure},
		{name: "empty password", hash: hash, errorWant: ErrInvalidInput},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &sessionCapture{err: tc.saveErr}
			s := NewServiceWithDependencies(loginUsers{credentials: Credentials{UserID: uuid.New(), PasswordHash: tc.hash}, err: tc.userErr}, store, passwords, brokenTokens{TokenManager: manager, generateErr: tc.generateErr, issueErr: tc.issueErr})
			result, err := s.Login(context.Background(), LoginInput{Username: "Alice", Password: tc.password})
			if !errors.Is(err, tc.errorWant) || result != (Tokens{}) || len(store.records) != tc.saves {
				t.Fatalf("result=%v err=%v saves=%d", result, err, len(store.records))
			}
		})
	}
	if _, err := NewService(nil).Login(context.Background(), LoginInput{}); !errors.Is(err, ErrAuthNotConfigured) {
		t.Fatal(err)
	}
}
