package auth

import (
	"context"
	"fmt"
)

func (s *Service) Register(ctx context.Context, input CreateInput) (Created, error) {
	if s.users == nil {
		return Created{}, ErrStoreNotConfigured
	}

	input = normalizeCreateInput(input)
	if err := validateCreateInput(input); err != nil {
		return Created{}, InvalidInputError{Err: err}
	}

	hash, err := s.passwords.Hash(input.Password)
	if err != nil {
		return Created{}, fmt.Errorf("hash password: %w", err)
	}

	id, err := s.users.CreateUser(ctx, CreateRecord{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hash,
	})
	if err != nil {
		return Created{}, err
	}

	return Created{ID: id}, nil
}
