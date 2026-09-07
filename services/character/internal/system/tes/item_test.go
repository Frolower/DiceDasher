package tes

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestItemFactoryVariants(t *testing.T) {
	for _, tt := range []struct {
		raw    string
		want   item
		absent string
	}{
		{`{"type":"gear","name":"Tool","code":"tool","bonus":1}`, equipment{}, "armor_level"},
		{`{"type":"weapon","name":"Gun","code":"gun","damage_kind":"physical","damage_value":2,"range_min":"short","range_max":"long","tags":["loud"]}`, weapon{}, "processor"},
		{`{"type":"armor","name":"Vest","code":"vest","armor_level":2,"agility_modifier":-1}`, armor{}, "damage_kind"},
		{`{"type":"neurocaster","name":"Deck","code":"deck","processor":1,"network":2,"graphics":3}`, neurocaster{}, "bonus"},
	} {
		t.Run(tt.want.kind(), func(t *testing.T) {
			var dto gearDTO
			if err := json.Unmarshal([]byte(tt.raw), &dto); err != nil {
				t.Fatal(err)
			}
			value, violations := newItem(dto)
			if len(violations) != 0 {
				t.Fatal(violations)
			}
			if reflect.TypeOf(value) != reflect.TypeOf(tt.want) {
				t.Fatalf("unexpected variant %T", value)
			}
			data, err := json.Marshal(dto)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err = json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			if _, ok := fields[tt.absent]; ok {
				t.Fatalf("foreign field persisted: %s", data)
			}
			var roundTrip gearDTO
			if err = json.Unmarshal(data, &roundTrip); err != nil {
				t.Fatal(err)
			}
			restored, errs := newItem(roundTrip)
			if len(errs) > 0 || !reflect.DeepEqual(value, restored) {
				t.Fatalf("round trip changed value: %s, %v", data, errs)
			}
		})
	}
}

func TestItemFactoryRejectsInvalidVariants(t *testing.T) {
	for _, raw := range []string{
		`{"type":"unknown","name":"X","code":"x"}`,
		`{"type":"gear","name":"X","code":"x","armor_level":1}`,
		`{"type":"armor","name":"X","code":"x","armor_level":1,"agility_modifier":-2,"bonus":1}`,
		`{"type":"weapon","name":"X","code":"x","damage_kind":"physical","range_min":"long","range_max":"short"}`,
		`{"type":"neurocaster","name":"X","code":"x","processor":0,"network":1,"graphics":1}`,
	} {
		var dto gearDTO
		if err := json.Unmarshal([]byte(raw), &dto); err != nil {
			t.Fatal(err)
		}
		value, violations := newItem(dto)
		if value != nil || len(violations) == 0 {
			t.Fatalf("accepted invalid item: %s", raw)
		}
	}
}

func TestItemFactoryAcceptsLegacyZeroFields(t *testing.T) {
	var dto gearDTO
	if err := json.Unmarshal([]byte(`{"type":"gear","name":"Tool","code":"tool","armor_level":0,"processor":0,"damage_kind":"","tags":[]}`), &dto); err != nil {
		t.Fatal(err)
	}
	if _, v := newItem(dto); len(v) > 0 {
		t.Fatal(v)
	}
}

func TestTypedItemsPersistInEveryInventory(t *testing.T) {
	raw := `{"user_id":"d7a92c4c-7c65-41d8-946a-4b94d3e721f9","type":"pc","character":{"name":"Mira","gear":[{"type":"gear","name":"Tool","code":"tool"}],"vehicle":{"SharedGear":[{"type":"gear","name":"Radio","code":"radio"}],"stats":{"gear":[{"type":"gear","name":"Kit","code":"kit"}]}}}}`
	created, status, err := (Character{}).CreateCharacter(context.Background(), json.RawMessage(raw))
	if err != nil || status != 201 {
		t.Fatalf("status %d: %v", status, err)
	}
	var stored struct {
		Character struct {
			Gear    []map[string]any
			Vehicle struct {
				SharedGear []map[string]any
				Stats      struct{ Gear []map[string]any }
			}
		}
	}
	if err = json.Unmarshal(created.Data, &stored); err != nil {
		t.Fatal(err)
	}
	for _, items := range [][]map[string]any{stored.Character.Gear, stored.Character.Vehicle.SharedGear, stored.Character.Vehicle.Stats.Gear} {
		if len(items) != 1 || items[0]["type"] != "gear" {
			t.Fatalf("missing item: %s", created.Data)
		}
		if _, ok := items[0]["armor_level"]; ok {
			t.Fatalf("foreign field: %s", created.Data)
		}
	}
}

func TestEveryItemIsValidatedWithoutCreationRules(t *testing.T) {
	r := validRequest()
	r.Rules = false
	r.CharacterList.Inventory = append(r.CharacterList.Inventory, gearDTO{Name: "Deck", Code: "deck", Type: neurocasterType, Processor: 1, Network: 1, Graphics: 0})
	if err := validateCreate(r); err == nil {
		t.Fatal("invalid second item accepted")
	}
	r.CharacterList.Inventory[1].Graphics = 1
	r.CharacterList.Inventory = append(r.CharacterList.Inventory, r.CharacterList.Inventory[1])
	if err := validateCreate(r); err != nil {
		t.Fatal(err)
	}
	r.Rules = true
	if err := validateCreate(r); err == nil {
		t.Fatal("multiple neurocasters accepted at creation")
	}
}
