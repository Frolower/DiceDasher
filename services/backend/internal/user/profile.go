package user

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (UserProfile, error) {
	if s.users == nil {
		return UserProfile{}, ErrStoreNotConfigured
	}
	if userID == uuid.Nil {
		return UserProfile{}, ErrUserNotFound
	}

	user, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		return UserProfile{}, err
	}

	return user, nil
}
