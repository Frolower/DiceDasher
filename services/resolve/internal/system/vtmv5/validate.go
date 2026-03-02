package vtmv5

import (
	"diceDasher/pkg/util"
	"errors"
)

func validateRoll(req rollRequest) error {
	var errs []error

	if req.Attribute < 1 {
		errs = append(errs, errors.New("attribute must be >= 1"))
	}
	if req.Skill < 0 {
		errs = append(errs, errors.New("skill must be >= 0"))
	}
	if req.Hunger < 0 {
		errs = append(errs, errors.New("hunger must be >= 0"))
	}
	if req.Target < 1 {
		errs = append(errs, errors.New("target must be >= 1"))
	}

	return errors.Join(errs...)
}

func validateReroll(req rerollRequest) error {
	var errs []error

	if len(req.MainRoll)+len(req.HungerRoll) == 0 {
		errs = append(errs, errors.New("no dice in a roll"))
	}
	if len(req.RerollIndex) == 0 {
		errs = append(errs, errors.New("no dice to reroll"))
	}
	if util.CountAbove(req.RerollIndex, len(req.MainRoll)-1) > 0 {
		errs = append(errs, errors.New("reroll index out of range"))
	}
	if req.Target < 1 {
		errs = append(errs, errors.New("target must be >= 1"))
	}

	return errors.Join(errs...)
}
