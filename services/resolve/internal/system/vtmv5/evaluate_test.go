package vtmv5

import (
	"diceDasher/pkg/dice"
	"encoding/json"
	"reflect"
	"testing"
)

func TestEvaluate(t *testing.T) {
	for _, tc := range []struct {
		name         string
		main, hunger []int
		target       int
		want         Outcome
	}{
		{"ordinary success", []int{5, 6, 9}, []int{2}, 2, Outcome{2, true, false, "none"}},
		{"ordinary failure", []int{5}, []int{2}, 1, Outcome{0, false, false, "none"}},
		{"bestial failure", []int{6}, []int{1}, 2, Outcome{1, false, true, "bestial failure"}},
		{"hunger one with success", []int{6}, []int{1}, 1, Outcome{1, true, false, "none"}},
		{"regular critical", []int{10, 10}, nil, 4, Outcome{4, true, true, "critical"}},
		{"mixed critical", []int{10}, []int{10}, 4, Outcome{4, true, true, "messy critical"}},
		{"hunger pair", nil, []int{10, 10}, 4, Outcome{4, true, true, "messy critical"}},
		{"unpaired ten", []int{10, 10}, []int{10}, 5, Outcome{5, true, true, "messy critical"}},
		{"multiple pairs", []int{10, 10, 10}, []int{10}, 8, Outcome{8, true, true, "messy critical"}},
		{"critical below target", []int{10, 10}, []int{1}, 5, Outcome{4, false, true, "critical"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			beforeMain := append([]int(nil), tc.main...)
			beforeHunger := append([]int(nil), tc.hunger...)
			if got := Evaluate(tc.main, tc.hunger, tc.target); got != tc.want {
				t.Fatalf("got %+v want %+v", got, tc.want)
			}
			if !reflect.DeepEqual(tc.main, beforeMain) || !reflect.DeepEqual(tc.hunger, beforeHunger) {
				t.Fatal("inputs mutated")
			}
		})
	}
}

func sequence(t *testing.T, values ...int) dice.Generator {
	t.Helper()
	i := 0
	return dice.NewGenerator(func(n int) int {
		if i >= len(values) {
			t.Fatal("unexpected random draw")
		}
		value := values[i]
		i++
		if value < 1 || value > n {
			t.Fatalf("invalid test die %d for d%d", value, n)
		}
		return value - 1
	})
}

func TestRollAndRerollEvaluateIdentically(t *testing.T) {
	r := Resolver{Dice: sequence(t, 10, 6, 10)}
	roll, err := r.resolveRoll(json.RawMessage(`{"attribute":3,"hunger":1,"target":5}`))
	if err != nil {
		t.Fatalf("roll: %v", err)
	}
	if roll.Successes != 5 || !roll.Success || roll.CritType != "messy critical" {
		t.Fatalf("wrong roll: %+v", roll)
	}
	source := rollState{MainRoll: []int{2, 6}, HungerRoll: []int{10}, Target: 5}
	r.Dice = sequence(t, 10)
	reroll, err := r.continueRoll(source, []int{0})
	if err != nil {
		t.Fatalf("reroll: %v", err)
	}
	if !reflect.DeepEqual(reroll.MainRoll, roll.MainRoll) || !reflect.DeepEqual(reroll.HungerRoll, roll.HungerRoll) ||
		reroll.Successes != roll.Successes || reroll.Success != roll.Success || reroll.IsCritical != roll.IsCritical || reroll.CritType != roll.CritType {
		t.Fatalf("inconsistent reroll: %+v", reroll)
	}
	if source.MainRoll[0] != 2 {
		t.Fatal("source mutated")
	}
}

func TestCheckUsesInjectedGenerator(t *testing.T) {
	for _, value := range []int{5, 6} {
		r := Resolver{Dice: sequence(t, value)}
		got, err := r.resolveCheck(nil)
		if err != nil || got.Result != value || got.Success != (value >= 6) {
			t.Fatalf("check: %+v %v", got, err)
		}
	}
}
