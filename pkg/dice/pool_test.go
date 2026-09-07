package dice

import "testing"

func TestRollDiceBounds(t *testing.T) {
	max := int(^uint(0) >> 1)
	for _, tc := range [][2]int{{-1, 6}, {MaxDice + 1, 6}, {max, 6}, {1, 1}, {1, MaxSides + 1}} {
		if _, err := RollDice(tc[0], tc[1]); err == nil {
			t.Fatalf("accepted %v", tc)
		}
	}
	for _, tc := range [][2]int{{0, 6}, {MaxDice, MaxSides}} {
		pool, err := NewPool(tc[0], tc[1])
		if err != nil {
			t.Fatal(err)
		}
		rolls, err := pool.Roll()
		if err != nil {
			t.Fatal(err)
		}
		if len(rolls) != tc[0] {
			t.Fatal("wrong count")
		}
		for _, v := range rolls {
			if v < 1 || v > tc[1] {
				t.Fatalf("invalid die: %d", v)
			}
		}
	}
	var zero Pool
	if _, err := zero.Roll(); err == nil {
		t.Fatal("zero-value pool must fail safely")
	}
}

func TestAddCountsOverflow(t *testing.T) {
	max := int(^uint(0) >> 1)
	min := -max - 1
	for _, values := range [][]int{{max, 1}, {min, -1}, {max, max, min}} {
		if _, err := AddCounts(values...); err == nil {
			t.Fatalf("accepted overflow: %v", values)
		}
	}
	for _, tc := range []struct {
		values []int
		want   int
	}{{[]int{max, -max}, 0}, {[]int{min, max}, -1}, {[]int{2, 3, -4}, 1}} {
		got, err := AddCounts(tc.values...)
		if err != nil || got != tc.want {
			t.Fatalf("%v: %d %v", tc.values, got, err)
		}
	}
}

func TestRerollLimits(t *testing.T) {
	previous := make([]int, MaxDice+1)
	if _, err := RerollKeepingValues(previous, []int{1, 6}, 6); err == nil {
		t.Fatal("oversized reroll accepted")
	}
	if _, err := RerollSpecificValues(previous, []int{0}, 6); err == nil {
		t.Fatal("oversized reroll accepted")
	}
	if _, err := RollDie(MaxSides + 1); err == nil {
		t.Fatal("oversized die accepted")
	}
}
