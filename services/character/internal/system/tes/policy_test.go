package tes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"diceDasher/services/character/internal/system"
	"github.com/google/uuid"
)

func validRequest() createRequest {
	item := gear{Name: "Tool", Code: "tool", Type: gearType}
	return createRequest{UserID: uuid.New(), Type: pc, Rules: true, CharacterList: characterList{
		Name: "Mira", Archetype: "artist", FavouriteSong: "Song", Stats: stats{4, 4, 4, 4},
		Derivatives: derivatives{4, 4}, Talents: []string{"athlete"}, Dream: "Dream", Flaw: "Flaw",
		Inventory: []gear{item}, Cash: 100, Journey: journey{"Goal", "Threat"},
		Vehicle: vehicle{VehicleType: "4wdCar", Model: "Car", Fuel: "gasoline", Stats: vehicleStats{Speed: 1, Hull: 1}, SharedGear: []gear{item, item, item}},
	}}
}

func TestCreationPolicies(t *testing.T) {
	for _, rules := range []bool{false, true} {
		t.Run(map[bool]string{false: "free", true: "strict"}[rules], func(t *testing.T) {
			r := validRequest()
			r.Rules = rules
			if err := validateCreate(r); err != nil {
				t.Fatal(err)
			}
			r.CharacterList.Name = "  "
			r.CharacterList.Inventory = append(r.CharacterList.Inventory, gear{Name: "Broken", Code: "broken", Type: gearType, Price: -1})
			err := validateCreate(r)
			var validation *system.ValidationError
			if !errors.As(err, &validation) {
				t.Fatalf("expected structured validation, got %v", err)
			}
			found := map[string]bool{}
			for _, v := range validation.Violations {
				found[v.Field+"/"+v.Code] = true
			}
			for _, key := range []string{"character.name/required", "character.gear[1].price/negative_value"} {
				if !found[key] {
					t.Errorf("missing %s: %v", key, err)
				}
			}
		})
	}
}

func TestFreeModeSkipsCreationLimits(t *testing.T) {
	r := validRequest()
	r.Rules = false
	r.CharacterList.Stats.Strength = 10
	r.CharacterList.Derivatives.Health = 99
	r.CharacterList.Cash = 99999
	r.CharacterList.Talents = nil
	r.CharacterList.Inventory = nil
	r.CharacterList.Vehicle = vehicle{}
	if err := validateCreate(r); err != nil {
		t.Fatal(err)
	}
	r.Rules = true
	if err := validateCreate(r); err == nil {
		t.Fatal("strict policy accepted invalid starting sheet")
	}
}

func TestCreateCharacterContract(t *testing.T) {
	for _, tt := range []struct {
		name, raw string
		status    int
	}{
		{"minimal", `{"user_id":"d7a92c4c-7c65-41d8-946a-4b94d3e721f9","type":"pc","character":{"name":"Mira"}}`, 201},
		{"missing sheet", `{"user_id":"d7a92c4c-7c65-41d8-946a-4b94d3e721f9","type":"pc"}`, 422},
		{"null sheet", `{"user_id":"d7a92c4c-7c65-41d8-946a-4b94d3e721f9","type":"pc","character":null}`, 422},
		{"missing owner", `{"type":"pc","character":{"name":"Mira"}}`, 400},
		{"invalid JSON", `{`, 400},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, status, err := (Character{}).CreateCharacter(context.Background(), json.RawMessage(tt.raw))
			if status != tt.status {
				t.Fatalf("status %d: %v", status, err)
			}
			if status == http.StatusCreated && err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIntegrityBoundaries(t *testing.T) {
	for _, tt := range []struct {
		name   string
		change func(*characterList)
	}{
		{"long name", func(c *characterList) { c.Name = strings.Repeat("я", 256) }},
		{"negative stat", func(c *characterList) { c.Stats.Wits = -1 }},
		{"negative bliss", func(c *characterList) { c.Bliss.Permanent = -1 }},
		{"unknown item", func(c *characterList) { c.Inventory[0].Type = "unknown" }},
		{"invalid weapon", func(c *characterList) { c.Inventory[0].Type = weaponType }},
		{"shared item", func(c *characterList) { c.Vehicle.SharedGear[0].Price = -1 }},
		{"vehicle item", func(c *characterList) { c.Vehicle.Stats.Gear = []gear{{Type: gearType}} }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := validRequest()
			r.Rules = false
			tt.change(&r.CharacterList)
			if validateCreate(r) == nil {
				t.Fatal("accepted invalid sheet")
			}
		})
	}
	r := validRequest()
	r.CharacterList.Name = strings.Repeat("я", 255)
	r.CharacterList.Inventory = []gear{{Name: "Armor", Code: "armor", Type: armorType, ArmorLevel: 1, AgilityModifier: -2}}
	if err := validateCreate(r); err != nil {
		t.Fatal(err)
	}
}
