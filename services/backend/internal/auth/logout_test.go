package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type revokeCapture struct {
	SessionStore
	ctx   context.Context
	hash  string
	now   time.Time
	calls int
	err   error
}

func (s *revokeCapture) RevokeSessionByRefreshHash(ctx context.Context, hash string, now time.Time) error {
	s.ctx, s.hash, s.now = ctx, hash, now
	s.calls++
	return s.err
}

func TestLogoutRevokesHashedRefreshToken(t *testing.T) {
	manager, err := NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
	if err != nil {
		t.Fatal(err)
	}
	failure := errors.New("database unavailable")
	for _, storeErr := range []error{nil, failure} {
		name := "success"
		if storeErr != nil {
			name = "repository error"
		}
		t.Run(name, func(t *testing.T) {
			store := &revokeCapture{err: storeErr}
			s := NewServiceWithDependencies(loginUsers{}, store, nil, manager)
			now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
			s.now = func() time.Time { return now }
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err := s.Logout(ctx, "secret-refresh")
			if !errors.Is(err, storeErr) {
				t.Fatalf("unexpected error: %v", err)
			}
			if store.calls != 1 || store.ctx != ctx || store.hash != manager.HashRefresh("secret-refresh") || store.hash == "secret-refresh" || !store.now.Equal(now) {
				t.Fatal("incorrect session revocation arguments")
			}
		})
	}
}

func TestLogoutUnconfigured(t *testing.T) {
	if err := NewService(nil).Logout(context.Background(), "refresh"); !errors.Is(err, ErrAuthNotConfigured) {
		t.Fatalf("expected configuration error, got %v", err)
	}
}
