package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

type authenticationStore struct {
	SessionStore
	session Session
	err     error
	calls   int
	ctx     context.Context
	id      uuid.UUID
}

func (s *authenticationStore) FindSessionByID(ctx context.Context, id uuid.UUID) (Session, error) {
	s.calls++
	s.ctx, s.id = ctx, id
	return s.session, s.err
}
func TestAuthenticate(t *testing.T) {
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	failure := errors.New("database unavailable")
	for _, tc := range []struct {
		name                                      string
		revoked, mismatch, invalidJWT, expiredJWT bool
		lifetime                                  time.Duration
		storeErr, wantErr                         error
	}{
		{name: "active session", lifetime: time.Hour},
		{name: "invalid JWT", invalidJWT: true, wantErr: ErrInvalidToken},
		{name: "expired JWT", expiredJWT: true, wantErr: ErrInvalidToken},
		{name: "missing session", storeErr: fmt.Errorf("find: %w", ErrSessionNotFound), wantErr: ErrSessionNotFound},
		{name: "database failure", storeErr: failure, wantErr: failure},
		{name: "revoked session", revoked: true, lifetime: time.Hour, wantErr: ErrSessionNotFound},
		{name: "expired session", lifetime: -time.Second, wantErr: ErrSessionNotFound},
		{name: "session expires exactly now", wantErr: ErrSessionNotFound},
		{name: "wrong user", mismatch: true, lifetime: time.Hour, wantErr: ErrSessionNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			manager, err := NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
			if err != nil {
				t.Fatal(err)
			}
			manager.now = func() time.Time { return now }
			userID, sessionID := uuid.New(), uuid.New()
			raw, err := manager.IssueAccess(userID, sessionID, now.Add(time.Minute))
			if err != nil {
				t.Fatal(err)
			}
			if tc.invalidJWT {
				raw = "invalid.jwt"
			}
			if tc.expiredJWT {
				manager.now = func() time.Time { return now.Add(time.Minute) }
			}
			store := &authenticationStore{session: Session{ID: sessionID, UserID: userID, ExpiresAt: now.Add(tc.lifetime)}, err: tc.storeErr}
			if tc.revoked {
				store.session.RevokedAt = &now
			}
			if tc.mismatch {
				store.session.UserID = uuid.New()
			}
			service := NewServiceWithDependencies(nil, store, nil, manager)
			service.now = func() time.Time { return now }
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			claims, err := service.Authenticate(ctx, raw)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error=%v want=%v", err, tc.wantErr)
			}
			if tc.wantErr != nil && claims != (AccessClaims{}) {
				t.Fatal("claims returned on failure")
			}
			if tc.invalidJWT || tc.expiredJWT {
				if store.calls != 0 {
					t.Fatal("invalid JWT queried database")
				}
				return
			}
			if store.calls != 1 || store.ctx != ctx || store.id != sessionID {
				t.Fatal("wrong session lookup")
			}
			if tc.wantErr == nil && (claims.UserID != userID || claims.SessionID != sessionID || !claims.ExpiresAt.Equal(now.Add(time.Minute))) {
				t.Fatalf("wrong claims: %v", claims)
			}
		})
	}
}
func TestAuthenticateUnconfigured(t *testing.T) {
	manager, err := NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
	if err != nil {
		t.Fatal(err)
	}
	for _, service := range []*Service{
		NewService(nil),
		NewServiceWithDependencies(nil, nil, nil, manager),
		NewServiceWithDependencies(nil, &authenticationStore{}, nil, nil),
	} {
		claims, err := service.Authenticate(context.Background(), "access-token")
		if !errors.Is(err, ErrAuthNotConfigured) || claims != (AccessClaims{}) {
			t.Fatalf("claims=%v err=%v", claims, err)
		}
	}
}
