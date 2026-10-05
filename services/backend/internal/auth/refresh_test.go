package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

type refreshStore struct {
	SessionStore
	session            Session
	findErr, rotateErr error
	finds, rotations   int
	oldHash, newHash   string
	rotationTime       time.Time
}

func (s *refreshStore) FindSessionByRefreshHash(ctx context.Context, hash string) (Session, error) {
	s.finds++
	if hash != s.oldHash {
		return Session{}, ErrSessionNotFound
	}
	return s.session, s.findErr
}
func (s *refreshStore) RotateRefreshToken(ctx context.Context, id uuid.UUID, oldHash, newHash string, now time.Time) error {
	s.rotations++
	if id != s.session.ID || oldHash != s.oldHash {
		return ErrSessionNotFound
	}
	s.newHash, s.rotationTime = newHash, now
	if s.rotateErr != nil {
		return s.rotateErr
	}
	s.oldHash = newHash
	return nil
}

func TestRefresh(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	failure := errors.New("internal failure")
	for _, tc := range []struct {
		name                                               string
		lifetime                                           time.Duration
		revoked, empty                                     bool
		findErr, generateErr, issueErr, rotateErr, wantErr error
		rotations                                          int
	}{
		{name: "success", lifetime: SessionLifetime, rotations: 1},
		{name: "session ends before access", lifetime: time.Minute, rotations: 1},
		{name: "empty token", empty: true, wantErr: ErrSessionNotFound},
		{name: "missing session", findErr: ErrSessionNotFound, wantErr: ErrSessionNotFound},
		{name: "read failure", findErr: failure, wantErr: failure},
		{name: "revoked", lifetime: SessionLifetime, revoked: true, wantErr: ErrSessionNotFound},
		{name: "expired", lifetime: -time.Second, wantErr: ErrSessionNotFound},
		{name: "expires exactly now", wantErr: ErrSessionNotFound},
		{name: "random failure", lifetime: SessionLifetime, generateErr: failure, wantErr: failure},
		{name: "signature failure", lifetime: SessionLifetime, issueErr: failure, wantErr: failure},
		{name: "rotation failure", lifetime: SessionLifetime, rotateErr: failure, wantErr: failure, rotations: 1},
		{name: "concurrent rotation lost", lifetime: SessionLifetime, rotateErr: ErrSessionNotFound, wantErr: ErrSessionNotFound, rotations: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			manager, err := NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
			if err != nil {
				t.Fatal(err)
			}
			manager.now = func() time.Time { return now }
			store := &refreshStore{session: Session{ID: uuid.New(), UserID: uuid.New(), ExpiresAt: now.Add(tc.lifetime)}, oldHash: manager.HashRefresh("old-refresh"), findErr: tc.findErr, rotateErr: tc.rotateErr}
			if tc.revoked {
				store.session.RevokedAt = &now
			}
			service := NewServiceWithDependencies(nil, store, nil, brokenTokens{TokenManager: manager, generateErr: tc.generateErr, issueErr: tc.issueErr})
			clockCalls := 0
			service.now = func() time.Time { clockCalls++; return now.Add(time.Duration(clockCalls-1) * time.Millisecond) }
			token := "old-refresh"
			if tc.empty {
				token = ""
			}
			result, err := service.Refresh(context.Background(), token)
			if !errors.Is(err, tc.wantErr) || store.rotations != tc.rotations {
				t.Fatalf("err=%v rotations=%d", err, store.rotations)
			}
			if tc.wantErr != nil {
				if result != (Tokens{}) {
					t.Fatal("returned tokens on failure")
				}
				if tc.empty && store.finds != 0 {
					t.Fatal("empty token queried repository")
				}
				return
			}
			claims, err := manager.VerifyAccess(result.AccessToken)
			expectedExpiry := now.Add(AccessLifetime)
			if store.session.ExpiresAt.Before(expectedExpiry) {
				expectedExpiry = store.session.ExpiresAt
			}
			if err != nil || claims.UserID != store.session.UserID || claims.SessionID != store.session.ID || !claims.ExpiresAt.Equal(expectedExpiry) {
				t.Fatalf("wrong JWT: %v %v", claims, err)
			}
			if result.RefreshToken == "" || result.RefreshToken == token || store.newHash != manager.HashRefresh(result.RefreshToken) || store.newHash == result.RefreshToken || !result.RefreshExpiresAt.Equal(store.session.ExpiresAt) || !result.AccessExpiresAt.Equal(expectedExpiry) {
				t.Fatal("wrong rotated tokens")
			}
			if !store.rotationTime.Equal(now.Add(time.Millisecond)) {
				t.Fatal("rotation did not use fresh time")
			}
			if reused, err := service.Refresh(context.Background(), token); !errors.Is(err, ErrSessionNotFound) || reused != (Tokens{}) {
				t.Fatal("old token was accepted twice")
			}
		})
	}
}
func TestRefreshUnconfigured(t *testing.T) {
	manager, err := NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []*Service{NewService(nil), NewServiceWithDependencies(nil, nil, nil, manager), NewServiceWithDependencies(nil, &refreshStore{}, nil, nil)} {
		if _, err := s.Refresh(context.Background(), "token"); !errors.Is(err, ErrAuthNotConfigured) {
			t.Fatalf("wrong error: %v", err)
		}
	}
}
