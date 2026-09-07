package system

import "errors"

type ErrorKind string

const (
	BadRequest ErrorKind = "bad_request"
	Validation ErrorKind = "validation"
)

// RequestError distinguishes invalid commands from execution failures.
type RequestError struct {
	Kind  ErrorKind
	Cause error
}

func (e *RequestError) Error() string           { return e.Cause.Error() }
func (e *RequestError) Unwrap() error           { return e.Cause }
func Invalid(kind ErrorKind, cause error) error { return &RequestError{Kind: kind, Cause: cause} }
func IsValidation(err error) bool {
	var e *RequestError
	return errors.As(err, &e) && e.Kind == Validation
}

var ErrUnknownSystem = errors.New("unknown system")
