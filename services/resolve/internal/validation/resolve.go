package validation

import (
	"diceDasher/services/resolve/internal/model"
	"errors"
)

func ValidateResolve(r model.Resolve) error {
	var errs []error

	if r.Number < 1 {
		errs = append(errs, errors.New("number of dice must be at least 1"))
	}
	if r.Size < 2 {
		errs = append(errs, errors.New("size of dice must be at least 2"))
	}

	return errors.Join(errs...)
}

func ValidateResolveTES(r model.ResolveTES) error {
	var errs []error

	if r.Attr < 1 {
		errs = append(errs, errors.New("attribute score must be at least 1"))
	}

	if r.Assist < 0 || r.Assist > 3 {
		errs = append(errs, errors.New("assist score must be between 0 and 3"))
	}

	if r.Gear < 0 {
		errs = append(errs, errors.New("gear score must be at least 0"))
	}

	if r.Attr+r.Assist+r.Gear+r.Modificator < 1 {
		errs = append(errs, errors.New("negative modificator score is greater than combined roll pull"))
	}

	if r.Target < 0 || r.Target > r.Attr+r.Assist+r.Gear {
		errs = append(errs, errors.New("target score must be between 0 and total dice number"))
	}

	return errors.Join(errs...)
}

//func ValidateResolveVtMV5 (r model.ResolveVtMV5) error {}
