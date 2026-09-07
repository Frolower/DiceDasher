package tes

import (
	"diceDasher/pkg/dice"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func rollPools(r rollRequest) (dice.Pool, dice.Pool, error) {
	fail := func(err error) (dice.Pool, dice.Pool, error) { return dice.Pool{}, dice.Pool{}, err }
	if r.Attr < 1 {
		return fail(errors.New("attribute score must be at least 1"))
	}
	if r.Assist < 0 || r.Assist > 3 {
		return fail(errors.New("assist score must be between 0 and 3"))
	}
	if r.Gear < 0 {
		return fail(errors.New("gear score must be at least 0"))
	}
	base, err := dice.AddCounts(r.Attr, r.Assist, r.Gear)
	if err != nil {
		return fail(err)
	}
	attributes, err := dice.AddCounts(r.Attr, r.Assist, r.Modificator)
	if err != nil {
		return fail(err)
	}
	attributePool, err := dice.NewPool(attributes, dieSize)
	if err != nil {
		return fail(fmt.Errorf("attribute pool: %w", err))
	}
	gearPool, err := dice.NewPool(r.Gear, dieSize)
	if err != nil {
		return fail(fmt.Errorf("gear pool: %w", err))
	}
	// Both components are bounded, so their sum cannot overflow.
	total := attributePool.Count() + gearPool.Count()
	if total < 1 || total > dice.MaxDice {
		return fail(fmt.Errorf("total dice count must be between 1 and %d", dice.MaxDice))
	}
	if r.Target < 0 || r.Target > base {
		return fail(errors.New("target score must be between 0 and unmodified total dice number"))
	}
	return attributePool, gearPool, nil
}

func validatePush(req pushRequest) error {
	var errs []error

	if req.RecordID == uuid.Nil {
		errs = append(errs, errors.New("record_id is required"))
	}

	return errors.Join(errs...)
}
