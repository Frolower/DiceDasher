package auth

import "errors"

var (
	ErrInvalidToken       = errors.New("invalid access token")
	ErrSessionNotFound    = errors.New("session not found")
	ErrAuthNotConfigured  = errors.New("authentication dependencies are not configured")
	ErrNotImplemented     = errors.New("authentication operation is not implemented")
	ErrUserNotFound       = errors.New("user not found")
	ErrUnableToLogin      = errors.New("invalid username or password")
	ErrAlreadyExists      = errors.New("username or email already exists")
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
