package dice

import (
	"fmt"
	"math/rand/v2"
)

// Generator owns its random source; the zero value uses the concurrency-safe
// package source. Injected sources must be safe for their intended usage.
type Generator struct{ intN func(int) int }

// NewGenerator accepts an IntN-compatible function returning values in [0, n).
func NewGenerator(intN func(int) int) Generator { return Generator{intN: intN} }

func (g Generator) next(sides int) (int, error) {
	source := g.intN
	if source == nil {
		source = rand.IntN
	}
	value := source(sides)
	if value < 0 || value >= sides {
		return 0, fmt.Errorf("random source returned an invalid value")
	}
	return value + 1, nil
}

func (g Generator) RollDie(sides int) (int, error) {
	if _, err := NewPool(1, sides); err != nil {
		return 0, err
	}
	return g.next(sides)
}

func (g Generator) RollDice(count int, sides int) ([]int, error) {
	if _, err := NewPool(count, sides); err != nil {
		return nil, err
	}

	rolls := make([]int, count)
	for i := 0; i < count; i++ {
		value, err := g.next(sides)
		if err != nil {
			return nil, err
		}
		rolls[i] = value
	}
	return rolls, nil
}

func (g Generator) RerollKeepingValues(previous []int, except []int, sides int) ([]int, error) {
	if _, err := NewPool(len(previous), sides); err != nil {
		return nil, err
	}

	exceptSet := make(map[int]struct{}, len(except))
	for _, v := range except {
		exceptSet[v] = struct{}{}
	}

	out := make([]int, len(previous))
	copy(out, previous)

	for i, v := range out {
		if _, keep := exceptSet[v]; keep {
			continue
		}
		d, err := g.RollDie(sides)
		if err != nil {
			return nil, err
		}
		out[i] = d
	}

	return out, nil
}

func (g Generator) RerollSpecificValues(previous []int, index []int, sides int) ([]int, error) {
	if _, err := NewPool(len(previous), sides); err != nil {
		return nil, err
	}

	indexSet := make(map[int]struct{}, len(index))
	for _, idx := range index {
		indexSet[idx] = struct{}{}
	}

	out := make([]int, len(previous))
	copy(out, previous)

	for i := range out {
		if _, ok := indexSet[i]; ok {
			d, err := g.RollDie(sides)
			if err != nil {
				return nil, err
			}
			out[i] = d
		}
	}

	return out, nil
}

func RollDie(sides int) (int, error) { return (Generator{}).RollDie(sides) }

func RollDice(count int, sides int) ([]int, error) { return (Generator{}).RollDice(count, sides) }

func RerollKeepingValues(previous []int, except []int, sides int) ([]int, error) {
	return (Generator{}).RerollKeepingValues(previous, except, sides)
}

func RerollSpecificValues(previous []int, index []int, sides int) ([]int, error) {
	return (Generator{}).RerollSpecificValues(previous, index, sides)
}
