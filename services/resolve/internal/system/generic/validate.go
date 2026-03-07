package generic

import (
	"errors"
)

func validate(r request) error {
	var errs []error

	if r.Number < 1 {
		errs = append(errs, errors.New("number of dice must be at least 1"))
	}
	if r.Size < 2 {
		errs = append(errs, errors.New("size of dice must be at least 2"))
	}

	return errors.Join(errs...)
}
