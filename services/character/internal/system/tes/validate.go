package tes

import (
	"errors"
	"fmt"
	"slices"
)

const (
	npc = "npc"
	pc  = "pc"

	lowerBound = 2
	upperBound = 6

	secondTalentBaseline = 15

	neurocasterType = "neurocaster"
	gearType        = "gear"
	weaponType      = "weapon"
	armorType       = "armor"

	toughTalent   = "tough"
	dreamerTalent = "dreamer"
)

func validateDerivatives(c characterList) error {
	var errs []error
	if slices.Contains(c.Talents, toughTalent) {
		if c.Derivatives.Health != (c.Stats.Strength+c.Stats.Agility+1)/2+2 {
			errs = append(errs, errors.New("health should be equal to (STR + AGL) / 2, rounded up, plus talent bonus"))
		}
	} else {
		if c.Derivatives.Health != (c.Stats.Strength+c.Stats.Agility+1)/2 {
			errs = append(errs, errors.New("health should be equal to (STR + AGL) / 2, rounded up, plus talent bonus"))
		}
	}

	if slices.Contains(c.Talents, dreamerTalent) {
		if c.Derivatives.Hope != (c.Stats.Empathy+c.Stats.Wits+1)/2+2 {
			errs = append(errs, errors.New("hope should be equal to (EMP + WIT) / 2, rounded up, plus talent bonus"))
		}
	} else {
		if c.Derivatives.Hope != (c.Stats.Empathy+c.Stats.Wits+1)/2 {
			errs = append(errs, errors.New("hope should be equal to (EMP + WIT) / 2, rounded up, plus talent bonus"))
		}
	}

	return errors.Join(errs...)
}

func validateArchetype(archetype string) error {
	var archetypes = map[string]struct{}{
		"artist":       {},
		"criminal":     {},
		"devotee":      {},
		"doctor":       {},
		"dronePilot":   {},
		"investigator": {},
		"outsider":     {},
		"runawayKid":   {},
		"scientist":    {},
		"veteran":      {}}

	if _, ok := archetypes[archetype]; !ok {
		return fmt.Errorf("%s is invalid archetype", archetype)
	}
	return nil
}

func (s stats) Validate() error {
	var errs []error

	checks := map[string]int{
		"strength": s.Strength,
		"agility":  s.Agility,
		"wits":     s.Wits,
		"empathy":  s.Empathy,
	}

	for name, v := range checks {
		if v < lowerBound || v > upperBound {
			errs = append(errs, fmt.Errorf("%s must be between 2 and 6, got %d", name, v))
		}
	}
	return errors.Join(errs...)
}

// Only the creation policy limits the number of neurocasters.
func validateInventory(inventory []gearDTO) error {
	count := 0
	for _, item := range inventory {
		if item.Type == neurocasterType {
			count++
		}
	}
	if count > 1 {
		return errors.New("maximum 1 neurocaster allowed during character creation")
	}
	return nil
}

func validateStartingCash(a string, m int) error {
	ranges := map[string]struct{ min, max int }{
		"artist":       {100, 600},
		"criminal":     {20, 120},
		"devotee":      {100, 600},
		"doctor":       {200, 1200},
		"dronePilot":   {0, 0},
		"investigator": {100, 600},
		"outsider":     {10, 60},
		"runawayKid":   {10, 60},
		"scientist":    {200, 1200},
		"veteran":      {20, 120},
	}

	r, ok := ranges[a]
	if !ok {
		return fmt.Errorf("unknown archetype: %s", a)
	}

	if m < r.min || m > r.max {
		return fmt.Errorf("%s is expected to have from %d to %d starting cash, got %d", a, r.min, r.max, m)
	}

	return nil
}

func validateJourney(j journey) error {
	var errs []error

	if j.Goal == "" {
		errs = append(errs, errors.New("goal is empty"))
	}
	if j.Threat == "" {
		errs = append(errs, errors.New("threat is empty"))
	}

	return errors.Join(errs...)
}

func validateTension(t []tension) error {
	var errs []error

	for _, t := range t {
		if t.TravellerName == "" {
			errs = append(errs, errors.New("traveller name is empty"))
		}
		if t.Tension < 1 {
			errs = append(errs, errors.New("tension should be greater or equal to 1"))
		}
		if t.Tension > 2 {
			errs = append(errs, errors.New("tension can't be greater than 2"))
		}
	}

	return errors.Join(errs...)
}

func validateTalents(talents []string, s int) error {
	var errs []error

	var validTalents = map[string]struct{}{
		"athlete":        {},
		"backstabber":    {},
		"biker":          {},
		"bladeFighter":   {},
		"boatman":        {},
		"bomber":         {},
		"bowman":         {},
		"charmer":        {},
		"clubFighter":    {},
		"conArtist":      {},
		"dataMiner":      {},
		"dirtyFighter":   {},
		"dramaQueen":     {},
		"dreamer":        {},
		"driver":         {},
		"droneOperator":  {},
		"electronics":    {},
		"evasive":        {},
		"gamer":          {},
		"hacker":         {},
		"hardened":       {},
		"intuition":      {},
		"leader":         {},
		"loneWolf":       {},
		"machinegunner":  {},
		"martialArtist":  {},
		"mechanic":       {},
		"medic":          {},
		"menacing":       {},
		"musician":       {},
		"neuroresistant": {},
		"nineLives":      {},
		"nurse":          {},
		"pilot":          {},
		"pistoleer":      {},
		"resilient":      {},
		"rider":          {},
		"scout":          {},
		"sleuth":         {},
		"sniper":         {},
		"speaker":        {},
		"stealthy":       {},
		"surgeon":        {},
		"technoBabbler":  {},
		"thief":          {},
		"tough":          {},
	}

	if s > secondTalentBaseline && len(talents) != 1 {
		errs = append(errs, errors.New("character must have 1 Talent"))
	} else if s <= secondTalentBaseline && len(talents) != 2 {
		errs = append(errs, errors.New("character must have 2 Talents"))
	}

	seen := map[string]struct{}{}

	for _, talent := range talents {
		if _, ok := validTalents[talent]; !ok {
			errs = append(errs, fmt.Errorf("%s is invalid talent", talent))
			continue
		}

		if _, ok := seen[talent]; ok {
			errs = append(errs, fmt.Errorf("%s talent is duplicated", talent))
			continue
		}

		seen[talent] = struct{}{}
	}

	return errors.Join(errs...)
}

func validateVehicle(v vehicle) error {
	var errs []error

	vehicleTypes := map[string]struct{}{
		"horse":                    {},
		"wagon":                    {},
		"bicycle":                  {},
		"motorcycle":               {},
		"dirtBike":                 {},
		"2wdCar":                   {},
		"4wdCar":                   {},
		"pickupTruck":              {},
		"van":                      {},
		"lightTruck":               {},
		"heavyTruck":               {},
		"bus":                      {},
		"rowboat":                  {},
		"smallSailingBoat":         {},
		"smallMotorBoat":           {},
		"helicopter":               {},
		"lightAirplane":            {},
		"smallCommercialDroneShip": {},
		"militaryDroneShip":        {},
	}

	errs = append(errs, validateVehicleStats(v.Stats))

	if _, ok := vehicleTypes[v.VehicleType]; !ok {
		errs = append(errs, errors.New("vehicle type is invalid"))
	}
	if v.Model == "" {
		errs = append(errs, errors.New("model is empty"))
	}
	if v.Passengers < 0 {
		errs = append(errs, errors.New("passengers should be greater or equal to 0"))
	}
	if v.Fuel == "" {
		errs = append(errs, errors.New("fuel field is empty"))
	}
	if len(v.SharedGear) != 3 {
		errs = append(errs, errors.New("you are required to have 3 shared gear elements"))
	}

	return errors.Join(errs...)
}

func validateVehicleStats(s vehicleStats) error {
	var errs []error

	if s.Maneuverability < 0 {
		errs = append(errs, errors.New("maneuverability should be greater or equal to 0"))
	}
	if s.Speed < 1 {
		errs = append(errs, errors.New("speed should be greater or equal to 1"))
	}
	if s.Hull < 1 {
		errs = append(errs, errors.New("hull should be greater or equal to 1"))
	}
	if s.Armor < 0 {
		errs = append(errs, errors.New("armor should be greater or equal to 0"))
	}

	for _, t := range s.Traits {
		errs = append(errs, validateCarTrait(t))
	}

	return errors.Join(errs...)
}

func validateCarTrait(t carTrait) error {
	var errs []error

	traits := map[string]struct{}{
		"maneuverability": {},
		"speed":           {},
		"hull":            {},
		"armor":           {},
	}

	if t.Name == "" {
		errs = append(errs, errors.New("car trait name is empty"))
	}
	if _, ok := traits[t.Name]; !ok {
		errs = append(errs, errors.New("car trait name is invalid"))
	}
	if t.Bonus < 1 {
		errs = append(errs, errors.New("bonus should be greater or equal to 1"))
	}

	return errors.Join(errs...)
}
