package user

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) UpdateUser(ctx context.Context, rec UserUpdateRecord) error {
	if s.users == nil {
		return ErrStoreNotConfigured
	}
	rec = normalizeUserUpdateInput(rec)
	err := validateUserUpdateRequest(rec)
	if err != nil {
		return InvalidInputError{Err: err}
	}

	err = s.users.UpdateUser(ctx, rec)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if s.users == nil {
		return ErrStoreNotConfigured
	}
	err := s.users.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
