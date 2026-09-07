package vtmv5

import (
	"diceDasher/pkg/dice"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func rollPools(req rollRequest) (dice.Pool, dice.Pool, error) {
	fail := func(err error) (dice.Pool, dice.Pool, error) { return dice.Pool{}, dice.Pool{}, err }
	if req.Attribute < 1 {
		return fail(errors.New("attribute must be >= 1"))
	}
	if req.Skill < 0 {
		return fail(errors.New("skill must be >= 0"))
	}
	if req.Hunger < 0 {
		return fail(errors.New("hunger must be >= 0"))
	}
	if req.Target < 1 {
		return fail(errors.New("target must be >= 1"))
	}
	total, err := dice.AddCounts(req.Attribute, req.Skill)
	if err != nil {
		return fail(err)
	}
	if _, err := dice.NewPool(total, dieSize); err != nil {
		return fail(err)
	}
	if req.Hunger > total {
		return fail(errors.New("hunger dice cannot exceed total dice count"))
	}
	main, err := dice.NewPool(total-req.Hunger, dieSize)
	if err != nil {
		return fail(err)
	}
	hunger, err := dice.NewPool(req.Hunger, dieSize)
	if err != nil {
		return fail(err)
	}
	return main, hunger, nil
}

func validateReroll(req rerollRequest) error {
	var errs []error

	if req.RecordID == uuid.Nil {
		errs = append(errs, errors.New("record_id is required"))
	}
	if len(req.RerollIndex) == 0 || len(req.RerollIndex) > dice.MaxDice {
		errs = append(errs, fmt.Errorf("reroll index count must be between 1 and %d", dice.MaxDice))
	}

	return errors.Join(errs...)
}

func validateRerollState(state rollState, indices []int) error {
	if len(indices) == 0 || len(indices) > len(state.MainRoll) {
		return errors.New("invalid number of reroll indices")
	}
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
