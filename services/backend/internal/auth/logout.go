package auth

import (
	"context"
)

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if s.users == nil || s.sessions == nil || s.passwords == nil || s.tokens == nil {
		return ErrAuthNotConfigured
	}

	hash := s.tokens.HashRefresh(refreshToken)

	err := s.sessions.RevokeSessionByRefreshHash(ctx, hash, s.now())
	if err != nil {
		return err
	}

	return nil
}
