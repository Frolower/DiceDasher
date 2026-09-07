package tes

import "errors"

func (n neurocaster) validate() error {
	var errs []error

	if n.Processor < 1 {
		errs = append(errs, errors.New("neurocaster processor should be greater or equal to 1"))
	}
	if n.Network < 1 {
		errs = append(errs, errors.New("neurocaster network should be greater or equal to 1"))
	}
	if n.Graphics < 1 {
		errs = append(errs, errors.New("neurocaster graphics should be greater or equal to 1"))
	}

	return errors.Join(errs...)
}

func (w weapon) validate() error {
	var errs []error

	if w.Bonus < 0 {
		errs = append(errs, errors.New("weapon bonus must be greater or equal to 0"))
	}
	if w.DamageKind == "explosive" {
		if w.BlastPower < 1 {
			errs = append(errs, errors.New("blast power should be greater or equal to 1"))
		}
	} else if w.DamageKind == "physical" {
		if w.DamageValue < 0 {
			errs = append(errs, errors.New("damage must be greater or equal to 0"))
		}
	} else {
		errs = append(errs, errors.New("unknown damage kind"))
	}

	errs = append(errs, validateRange(w.RangeMin, w.RangeMax))

	return errors.Join(errs...)
}

func validateRange(min, max string) error {
	var errs []error

	ranges := map[string]int{
		"engaged": 0,
		"short":   1,
		"medium":  2,
		"long":    3,
		"extreme": 4,
	}

	minr, minOk := ranges[min]
	if !minOk {
		errs = append(errs, errors.New("min range is not valid"))
	}

	maxr, maxOk := ranges[max]
	if !maxOk {
		errs = append(errs, errors.New("max range is not valid"))
	}

	if minOk && maxOk && minr > maxr {
		errs = append(errs, errors.New("min range is greater than max range"))
	}

	return errors.Join(errs...)
}

func (a armor) validate() error {
	var errs []error

	if a.AgilityModifier < -3 || a.AgilityModifier > -1 {
		errs = append(errs, errors.New("agility modifier should be between -1 and -3"))
	}
	if a.ArmorLevel < 1 {
		errs = append(errs, errors.New("armor level should be greater or equal to 1"))
	}

	return errors.Join(errs...)
}

func (g equipment) validate() error {
	var errs []error

	if g.Bonus < 0 {
		errs = append(errs, errors.New("gear bonus should be greater or equal to 0"))
	}

	return errors.Join(errs...)
}
