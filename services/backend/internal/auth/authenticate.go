package auth

import (
	"context"
)

func (s *Service) Authenticate(ctx context.Context, accessToken string) (AccessClaims, error) {
	if s.tokens == nil || s.sessions == nil {
		return AccessClaims{}, ErrAuthNotConfigured
	}

	claims, err := s.tokens.VerifyAccess(accessToken)
	if err != nil {
		return AccessClaims{}, err
	}

	session, err := s.sessions.FindSessionByID(ctx, claims.SessionID)
	if err != nil {
		return AccessClaims{}, err
	}

	now := s.now().UTC()
	if session.RevokedAt != nil || !session.ExpiresAt.After(now) || session.UserID != claims.UserID {
		return AccessClaims{}, ErrSessionNotFound
	}

	return claims, nil
}
