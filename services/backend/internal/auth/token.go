package auth

import (
	"github.com/google/uuid"
	"time"
)

// AccessClaims contains verified JWT claims rather than an unverified decoded payload.
type AccessClaims struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	ExpiresAt time.Time
}

// TokenManager decouples auth flows from the implementation of JWTs and random refresh tokens.
type TokenManager interface {
	IssueAccess(userID, sessionID uuid.UUID, expiresAt time.Time) (string, error)
	VerifyAccess(token string) (AccessClaims, error)
	GenerateRefresh() (string, error)
	HashRefresh(token string) string
}
