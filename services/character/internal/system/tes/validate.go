package tes

import (
	"errors"
	"fmt"
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
)

func validateCreate(r createRequest) error {
	var errs []error

	if r.Type != npc && r.Type != pc {
		errs = append(errs, fmt.Errorf("type must be pc or npc"))
	}

	if r.Rules {
		errs = append(errs, r.CharacterList.Validate())
	}

	return errors.Join(errs...)
}

func (c characterList) Validate() error {
	var errs []error
	statSum := c.Stats.Agility + c.Stats.Empathy + c.Stats.Strength + c.Stats.Wits

	if c.Name == "" {
		errs = append(errs, errors.New("character's name is empty"))
	}

	errs = append(errs, validateArchetype(c.Archetype))

	errs = append(errs, c.Stats.Validate())

	if c.FavouriteSong == "" {
		errs = append(errs, errors.New("favourite song is empty"))
	}

	if c.Derivatives.Health != (c.Stats.Strength+c.Stats.Agility+1)/2 {
		errs = append(errs, errors.New("health should be equal to STR + AGL / 2, rounded up"))
	}

	if c.Derivatives.Hope != (c.Stats.Empathy+c.Stats.Wits+1)/2 {
		errs = append(errs, errors.New("health should be equal to EMP + WIT / 2, rounded up"))
	}

	if c.Bliss.Bliss < 0 {
		errs = append(errs, errors.New("bliss should be greater or equal to 0"))
	}

	if c.Bliss.Permanent < 0 {
		errs = append(errs, errors.New("bliss should be greater or equal to 0"))
	}

	if statSum > secondTalentBaseline && len(c.Talents) != 1 {
		errs = append(errs, errors.New("character must have 1 Talent"))
	} else if statSum <= secondTalentBaseline && len(c.Talents) != 2 {
		errs = append(errs, errors.New("character must have 2 Talents"))
	}

	if c.Dream == "" {
		errs = append(errs, errors.New("dream is empty"))
	}

	if c.Flaw == "" {
		errs = append(errs, errors.New("flaw is empty"))
	}

	if len(c.Inventory) < 1 || len(c.Inventory) > 4 {
		errs = append(errs, errors.New("expected 0 or 1 neurocaster and 1 other gear element at the character creation"))
	}

	errs = append(errs, validateInventory(c.Inventory))

	errs = append(errs, validateStartingCash(c.Archetype, c.Cash))

	errs = append(errs, validateJorney(c.Journey))

	errs = append(errs, validateTension(c.Tension))

	errs = append(errs, validateVehicle(c.Vehicle))

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

func validateInventory(inventory []gear) error {
	var errs []error
	var n gear
	var g gear
	neurocasterCount := 0

	for _, g := range inventory {
		if g.Type == neurocasterType {
			neurocasterCount++
			n = g
		} else {
			g = g
		}
	}

	if neurocasterCount > 1 {
		errs = append(errs, errors.New("maximum 1 neurocaster allowed during character creation"))
	}

	if neurocasterCount == 1 {
		errs = append(errs, validateNeurocaster(n))
	}

	if g.Type == weaponType {
		errs = append(errs, validateWeapon(g))
	} else if g.Type == armorType {
		errs = append(errs, validateArmor(g))
	} else if g.Type == gearType {
		errs = append(errs, validateGear(g))
	}

	return errors.Join(errs...)
}

func validateNeurocaster(n gear) error {
	var errs []error

	if n.Name == "" {
		errs = append(errs, errors.New("neurocaster name is empty"))
	}
	if n.Code == "" {
		errs = append(errs, errors.New("neurocaster code is empty"))
	}
	if n.Price < 0 {
		errs = append(errs, errors.New("neurocaster price should be greater or equal to 0"))
	}
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

func validateWeapon(w gear) error {
	var errs []error

	if w.Name == "" {
		errs = append(errs, errors.New("weapon name is empty"))
	}
	if w.Code == "" {
		errs = append(errs, errors.New("weapon code is empty"))
	}
	if w.Price < 0 {
		errs = append(errs, errors.New("weapon price must be greater or equal to 0"))
	}
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

func validateArmor(a gear) error {
	var errs []error

	if a.Name == "" {
		errs = append(errs, errors.New("armor name is empty"))
	}
	if a.Code == "" {
		errs = append(errs, errors.New("armor code is empty"))
	}
	if a.Price < 0 {
		errs = append(errs, errors.New("armor price should be greater or equal to 0"))
	}
	if a.AgilityModifier < 0 {
		errs = append(errs, errors.New("agility modifier should be greater or equal to 0"))
	}
	if a.ArmorLevel < 1 {
		errs = append(errs, errors.New("armor level should be greater or equal to 1"))
	}

	return errors.Join(errs...)
}

func validateGear(g gear) error {
	var errs []error

	if g.Name == "" {
		errs = append(errs, errors.New("gear name is empty"))
	}
	if g.Code == "" {
		errs = append(errs, errors.New("gear code is empty"))
	}
	if g.Price < 0 {
		errs = append(errs, errors.New("gear price should be greater or equal to 0"))
	}
	if g.Bonus < 0 {
		errs = append(errs, errors.New("gear bonus should be greater or equal to 0"))
	}

	return errors.Join(errs...)
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

func validateJorney(j journey) error {
	var errs []error

	if len(j.Goal) == 0 {
		errs = append(errs, errors.New("goal is required"))
	}
	if len(j.Threat) == 0 {
		errs = append(errs, errors.New("threat is required"))
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

	return errors.Join(errs...)
}

func validateVehcleStats(s vehicleStats) error {
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
