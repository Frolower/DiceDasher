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
	id, err := uuid.NewV7()
	if err != nil {
		return Tokens{}, err
	}
	session := Session{
		ID: id,
		UserID: userID,
		CreatedAt: now,
		ExpiresAt: now.Add(SessionLifetime),
		UserAgent: userAgent,
	}
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
	// Preparing the JWT before INSERT avoids leaving a session behind if signing fails.
	// The client receives tokens only after the session is successfully saved.
	if err := s.sessions.CreateSession(ctx, session); err != nil {
		return Tokens{}, fmt.Errorf("save session: %w", err)
	}
	return Tokens{AccessToken: access, AccessExpiresAt: accessExpires, RefreshToken: refresh, RefreshExpiresAt: session.ExpiresAt}, nil
}
