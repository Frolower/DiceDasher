package auth

import "context"

func (s *Service) Refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	return Tokens{}, ErrNotImplemented
}
