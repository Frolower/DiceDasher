package auth

import (
	"context"
)

func (s *Service) Refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	if s.sessions == nil || s.tokens == nil {
		return Tokens{}, ErrAuthNotConfigured
	}

	if refreshToken == "" {
		return Tokens{}, ErrSessionNotFound
	}

	oldHash := s.tokens.HashRefresh(refreshToken)
	session, err := s.sessions.FindSessionByRefreshHash(ctx, oldHash)
	if err != nil {
		return Tokens{}, err
	}

	now := s.now().UTC()
	if session.RevokedAt != nil {
		return Tokens{}, ErrSessionNotFound
	}
	if !session.ExpiresAt.After(now) {
		return Tokens{}, ErrSessionNotFound
	}

	newRefresh, err := s.tokens.GenerateRefresh()
	if err != nil {
		return Tokens{}, err
	}

	newHash := s.tokens.HashRefresh(newRefresh)

	accessExpire := now.Add(AccessLifetime)
	if session.ExpiresAt.Before(accessExpire) {
		accessExpire = session.ExpiresAt
	}
	access, err := s.tokens.IssueAccess(session.UserID, session.ID, accessExpire)
	if err != nil {
		return Tokens{}, err
	}

	err = s.sessions.RotateRefreshToken(ctx, session.ID, oldHash, newHash, s.now().UTC())
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:      access,
		AccessExpiresAt:  accessExpire,
		RefreshToken:     newRefresh,
		RefreshExpiresAt: session.ExpiresAt,
	}, nil
}
