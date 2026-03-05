package tes

import (
	"errors"

	"github.com/google/uuid"
)

func validateRoll(r rollRequest) error {
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

func validatePush(req pushRequest) error {
	var errs []error

	if req.RecordID == uuid.Nil {
		errs = append(errs, errors.New("record_id is required"))
	}

	return errors.Join(errs...)
}

func validatePushRecord(req pushRecord) error {
	var errs []error

	if len(req.AttributeRolls) < 1 {
		errs = append(errs, errors.New("attribute rolls must be at least 1"))
	}
	if req.Target < 1 || req.Target > len(req.AttributeRolls) {
		errs = append(errs, errors.New("target score must be between 1 and total dice number"))
	}

	return errors.Join(errs...)
}
