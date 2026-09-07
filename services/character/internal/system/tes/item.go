package tes

import (
	"encoding/json"
	"strings"

	"diceDasher/services/character/internal/system"
)

// item is a closed set of domain variants. The transport DTO is converted by newItem.
// The discriminator comes from the concrete type, never from mutable item data.
type item interface {
	kind() string
	validate() error
}

type itemBase struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Price int    `json:"price"`
	Notes string `json:"notes"`
}

type equipment struct {
	itemBase
	Bonus int `json:"bonus"`
}

func (equipment) kind() string { return gearType }

type weapon struct {
	itemBase
	Bonus       int      `json:"bonus"`
	DamageValue int      `json:"damage_value"`
	DamageKind  string   `json:"damage_kind"`
	RangeMin    string   `json:"range_min"`
	RangeMax    string   `json:"range_max"`
	Tags        []string `json:"tags"`
	BlastPower  int      `json:"blast_power"`
}

func (weapon) kind() string { return weaponType }

type armor struct {
	itemBase
	ArmorLevel      int `json:"armor_level"`
	AgilityModifier int `json:"agility_modifier"`
}

func (armor) kind() string { return armorType }

type neurocaster struct {
	itemBase
	Processor int `json:"processor"`
	Network   int `json:"network"`
	Graphics  int `json:"graphics"`
}

func (neurocaster) kind() string { return neurocasterType }

// newItem returns no domain value if any common or variant-specific constraint fails.
// Violation paths are relative to the item so every inventory can use this factory.
func newItem(dto gearDTO) (item, []system.Violation) {
	var violations []system.Violation
	add := func(field, code, message string) {
		violations = append(violations, system.Violation{Field: field, Code: code, Message: message})
	}
	if strings.TrimSpace(dto.Name) == "" {
		add("name", "required", "name is required")
	}
	if strings.TrimSpace(dto.Code) == "" {
		add("code", "required", "code is required")
	}
	if dto.Price < 0 {
		add("price", "negative_value", "must be greater or equal to 0")
	}
	base := itemBase{Name: dto.Name, Code: dto.Code, Price: dto.Price, Notes: dto.Notes}
	var value item
	switch dto.Type {
	case gearType:
		value = equipment{itemBase: base, Bonus: dto.Bonus}
	case weaponType:
		value = weapon{itemBase: base, Bonus: dto.Bonus, DamageValue: dto.DamageValue, DamageKind: dto.DamageKind, RangeMin: dto.RangeMin, RangeMax: dto.RangeMax, Tags: append([]string(nil), dto.Tags...), BlastPower: dto.BlastPower}
	case armorType:
		value = armor{itemBase: base, ArmorLevel: dto.ArmorLevel, AgilityModifier: dto.AgilityModifier}
	case neurocasterType:
		value = neurocaster{itemBase: base, Processor: dto.Processor, Network: dto.Network, Graphics: dto.Graphics}
	default:
		add("type", "invalid_value", "unknown gear type")
	}
	if value != nil {
		for _, field := range []struct {
			name             string
			present, allowed bool
		}{
			{"bonus", dto.Bonus != 0, dto.Type == gearType || dto.Type == weaponType},
			{"damage_value", dto.DamageValue != 0, dto.Type == weaponType},
			{"damage_kind", dto.DamageKind != "", dto.Type == weaponType},
			{"range_min", dto.RangeMin != "", dto.Type == weaponType},
			{"range_max", dto.RangeMax != "", dto.Type == weaponType},
			{"tags", len(dto.Tags) != 0, dto.Type == weaponType},
			{"blast_power", dto.BlastPower != 0, dto.Type == weaponType},
			{"armor_level", dto.ArmorLevel != 0, dto.Type == armorType},
			{"agility_modifier", dto.AgilityModifier != 0, dto.Type == armorType},
			{"processor", dto.Processor != 0, dto.Type == neurocasterType},
			{"network", dto.Network != 0, dto.Type == neurocasterType},
			{"graphics", dto.Graphics != 0, dto.Type == neurocasterType},
		} {
			if field.present && !field.allowed {
				add(field.name, "incompatible_field", "field does not belong to item type "+dto.Type)
			}
		}
		if err := value.validate(); err != nil {
			add("", "invalid_item", err.Error())
		}
	}
	if len(violations) > 0 {
		return nil, violations
	}
	return value, nil
}

// MarshalJSON preserves the flat wire format while persisting only fields belonging
// to the selected variant. Legacy DTOs containing irrelevant zero fields remain valid.
func (dto gearDTO) MarshalJSON() ([]byte, error) {
	value, violations := newItem(dto)
	if len(violations) > 0 {
		return nil, &system.ValidationError{Violations: violations}
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err = json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	fields["type"], err = json.Marshal(value.kind())
	if err != nil {
		return nil, err
	}
	return json.Marshal(fields)
}
