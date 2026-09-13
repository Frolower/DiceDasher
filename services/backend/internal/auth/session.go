package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const AccessLifetime = 10 * time.Minute
const SessionLifetime = 30 * 24 * time.Hour

func (s *Service) createSession(ctx context.Context, userID uuid.UUID, userAgent string) (Tokens, error) {
	now := s.now().UTC().Truncate(time.Second)
	session := Session{ID: uuid.New(), UserID: userID, CreatedAt: now, ExpiresAt: now.Add(SessionLifetime), UserAgent: userAgent}
	refresh, err := s.tokens.GenerateRefresh()
	if err != nil {
		return Tokens{}, fmt.Errorf("generate refresh token: %w", err)
	}
	session.RefreshTokenHash = s.tokens.HashRefresh(refresh)
	accessExpires := now.Add(AccessLifetime)
	access, err := s.tokens.IssueAccess(userID, session.ID, accessExpires)
	if err != nil {
		return Tokens{}, fmt.Errorf("issue access token: %w", err)
	}
	// Подготовка JWT до INSERT не оставляет сессию при ошибке подписи.
	// Но клиент получает токены только после успешного сохранения сессии.
	if err := s.sessions.CreateSession(ctx, session); err != nil {
		return Tokens{}, fmt.Errorf("save session: %w", err)
	}
	return Tokens{AccessToken: access, AccessExpiresAt: accessExpires, RefreshToken: refresh, RefreshExpiresAt: session.ExpiresAt}, nil
}
