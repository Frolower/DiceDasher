package auth

import (
	"github.com/google/uuid"
	"time"
)

// AccessClaims — результат проверки JWT, а не произвольно декодированное тело.
type AccessClaims struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	ExpiresAt time.Time
}

// TokenManager отделяет сценарии auth от реализации JWT и случайных refresh-токенов.
type TokenManager interface {
	IssueAccess(userID, sessionID uuid.UUID, expiresAt time.Time) (string, error)
	VerifyAccess(token string) (AccessClaims, error)
	GenerateRefresh() (string, error)
	HashRefresh(token string) string
}
