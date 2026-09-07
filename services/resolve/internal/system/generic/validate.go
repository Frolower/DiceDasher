package generic

import (
	"diceDasher/pkg/dice"
	"errors"
)

func validate(r request) error {
	if r.Number < 1 {
		return errors.New("number of dice must be at least 1")
	}
	_, err := dice.NewPool(r.Number, r.Size)
	return err
}
