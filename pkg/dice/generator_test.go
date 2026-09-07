package dice

import (
	"reflect"
	"testing"
)

func TestGeneratorUsedByAllOperations(t *testing.T) {
	calls := 0
	g := NewGenerator(func(n int) int { calls++; return n - 1 })
	rolls, err := g.RollDice(2, 6)
	if err != nil || !reflect.DeepEqual(rolls, []int{6, 6}) {
		t.Fatalf("roll: %v %v", rolls, err)
	}
	previous := []int{1, 2, 6, 3}
	kept, err := g.RerollKeepingValues(previous, []int{1, 6}, 6)
	if err != nil || !reflect.DeepEqual(kept, []int{1, 6, 6, 6}) {
		t.Fatalf("keep: %v %v", kept, err)
	}
	selected, err := g.RerollSpecificValues(previous, []int{1}, 6)
	if err != nil || !reflect.DeepEqual(selected, []int{1, 6, 6, 3}) {
		t.Fatalf("specific: %v %v", selected, err)
	}
	if calls != 5 || !reflect.DeepEqual(previous, []int{1, 2, 6, 3}) {
		t.Fatal("extra draw or source mutation")
	}
	if _, err := g.RollDice(-1, 6); err == nil || calls != 5 {
		t.Fatal("invalid pool reached source")
	}
}
func TestGeneratorRejectsInvalidSource(t *testing.T) {
	for _, value := range []int{-1, 6} {
		g := NewGenerator(func(int) int { return value })
		if _, err := g.RollDice(1, 6); err == nil {
			t.Fatal("invalid random value accepted")
		}
	}
}
