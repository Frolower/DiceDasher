package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidInput       = errors.New("invalid user input")
	ErrNotImplemented     = errors.New("not implemented")
	ErrStoreNotConfigured = errors.New("user store is not configured")
	ErrAlreadyExists      = errors.New("user already exists")
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
