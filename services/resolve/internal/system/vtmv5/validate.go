package vtmv5

import (
	"errors"

	"github.com/google/uuid"
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

	if req.RecordID == uuid.Nil {
		errs = append(errs, errors.New("record_id is required"))
	}
	if len(req.RerollIndex) == 0 {
		errs = append(errs, errors.New("no dice to reroll"))
	}

	return errors.Join(errs...)
}

func validateRerollState(state rollState, indices []int) error {
	seen := make(map[int]bool, len(indices))
	for _, index := range indices {
		if index < 0 || index >= len(state.MainRoll) {
			return errors.New("reroll index out of range")
		}
		if seen[index] {
			return errors.New("duplicate reroll index")
		}
		seen[index] = true
	}
	return nil
}
