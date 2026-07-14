package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput       = errors.New("invalid user input")
	ErrStoreNotConfigured = errors.New("user store is not configured")
)

type InvalidInputError struct {
	Err error
}

func (e InvalidInputError) Error() string {
	return e.Err.Error()
}

func (e InvalidInputError) Unwrap() error {
	return e.Err
}

func (e InvalidInputError) Is(target error) bool {
	return target == ErrInvalidInput
}

type Store interface {
	CreateUser(ctx context.Context, rec CreateRecord) (uuid.UUID, error)
}

type Service struct {
	store  Store
	hasher PasswordHasher
}

func NewService(store Store) *Service {
	return &Service{
		store:  store,
		hasher: NewBcryptHasher(),
	}
}

func NewServiceWithHasher(store Store, hasher PasswordHasher) *Service {
	if hasher == nil {
		hasher = NewBcryptHasher()
	}
	return &Service{store: store, hasher: hasher}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Created, error) {
	if s.store == nil {
		return Created{}, ErrStoreNotConfigured
	}

	input = normalizeCreateInput(input)
	if err := validateCreateInput(input); err != nil {
		return Created{}, InvalidInputError{Err: err}
	}

	hash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Created{}, err
	}

	id, err := s.store.CreateUser(ctx, CreateRecord{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hash,
	})
	if err != nil {
		return Created{}, err
	}

	return Created{ID: id}, nil
}
