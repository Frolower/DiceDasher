package tes

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"diceDasher/services/character/internal/system"
)

// A policy appends violations without stopping at the first invalid field.
type validationPolicy func(characterList, *system.ValidationError)

func validateCreate(r createRequest) error {
	result := &system.ValidationError{}
	if r.Type != pc && r.Type != npc {
		result.Violations = append(result.Violations, system.Violation{Field: "type", Code: "invalid_value", Message: "type must be pc or npc"})
	}
	policies := []validationPolicy{validateIntegrity}
	if r.Rules {
		policies = append(policies, creationPolicy)
	}
	for _, policy := range policies {
		policy(r.CharacterList, result)
	}
	if len(result.Violations) == 0 {
		return nil
	}
	return result
}

func addViolation(result *system.ValidationError, field, code, message string) {
	result.Violations = append(result.Violations, system.Violation{Field: strings.TrimSuffix("character."+field, "."), Code: code, Message: message})
}

func validateIntegrity(c characterList, result *system.ValidationError) {
	if strings.TrimSpace(c.Name) == "" {
		addViolation(result, "name", "required", "name is required")
	}
	if utf8.RuneCountInString(c.Name) > 255 {
		addViolation(result, "name", "too_long", "name must contain at most 255 characters")
	}
	for _, field := range []struct {
		name  string
		value int
	}{
		{"stats.strength", c.Stats.Strength}, {"stats.agility", c.Stats.Agility},
		{"stats.wits", c.Stats.Wits}, {"stats.empathy", c.Stats.Empathy},
		{"derivatives.health", c.Derivatives.Health}, {"derivatives.hope", c.Derivatives.Hope},
		{"bliss.bliss", c.Bliss.Bliss}, {"bliss.permanent", c.Bliss.Permanent}, {"cash", c.Cash},
	} {
		if field.value < 0 {
			addViolation(result, field.name, "negative_value", "must be greater or equal to 0")
		}
	}
	validateItems(c.Inventory, "gear", result)
	validateItems(c.Vehicle.SharedGear, "vehicle.SharedGear", result)
	validateItems(c.Vehicle.Stats.Gear, "vehicle.stats.gear", result)
}

func validateItems(items []gearDTO, path string, result *system.ValidationError) {
	for i, item := range items {
		field := fmt.Sprintf("%s[%d]", path, i)
		_, violations := newItem(item)
		for _, violation := range violations {
			itemField := field
			if violation.Field != "" {
				itemField += "." + violation.Field
			}
			addViolation(result, itemField, violation.Code, violation.Message)
		}
	}
}

func creationPolicy(c characterList, result *system.ValidationError) {
	for _, field := range []struct{ name, value string }{
		{"favourite_song", c.FavouriteSong}, {"dream", c.Dream}, {"flaw", c.Flaw},
	} {
		if strings.TrimSpace(field.value) == "" {
			addViolation(result, field.name, "required", "required by creation rules")
		}
	}
	checks := []struct {
		field string
		err   error
	}{
		{"archetype", validateArchetype(c.Archetype)},
		{"stats", c.Stats.Validate()},
		{"derivatives", validateDerivatives(c)},
		{"talents", validateTalents(c.Talents, c.Stats.Strength+c.Stats.Agility+c.Stats.Wits+c.Stats.Empathy)},
		{"gear", validateInventory(c.Inventory)},
		{"cash", validateStartingCash(c.Archetype, c.Cash)},
		{"journey", validateJourney(c.Journey)},
		{"tension", validateTension(c.Tension)},
		{"vehicle", validateVehicle(c.Vehicle)},
	}
	for _, check := range checks {
		appendRuleViolations(result, check.field, check.err)
	}
	if len(c.Inventory) < 1 || len(c.Inventory) > 4 {
		addViolation(result, "gear", "creation_rule", "inventory must contain between 1 and 4 items during character creation")
	}
}

func appendRuleViolations(result *system.ValidationError, field string, err error) {
	if err == nil {
		return
	}
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for _, child := range joined.Unwrap() {
			appendRuleViolations(result, field, child)
		}
		return
	}
	addViolation(result, field, "creation_rule", err.Error())
}
