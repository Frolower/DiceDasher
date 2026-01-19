package validation

import (
	"diceDasher/services/resolve/internal/model"
	"errors"
)

func ValidateResolve(r model.Resolve) error {
	if r.Number < 1 {
		return errors.New("number of dice must be at least 1")
	}

	if r.Size < 2 {
		return errors.New("size of dice must be at least 2")
	}

	return nil
}
