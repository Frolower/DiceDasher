package dice

import (
	"math/rand/v2"
)

func RollDie(sides int) (int, error) {
	if _, err := NewPool(1, sides); err != nil {
		return 0, err
	}
	roll := rand.IntN(sides) + 1
	return roll, nil
}

func RollDice(count int, sides int) ([]int, error) {
	if _, err := NewPool(count, sides); err != nil {
		return nil, err
	}

	rolls := make([]int, count)
	for i := 0; i < count; i++ {
		rolls[i] = rand.IntN(sides) + 1
	}
	return rolls, nil
}

func RerollKeepingValues(previous []int, except []int, sides int) ([]int, error) {
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
		d, err := RollDie(sides)
		if err != nil {
			return nil, err
		}
		out[i] = d
	}

	return out, nil
}

func RerollSpecificValues(previous []int, index []int, sides int) ([]int, error) {
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
			d, err := RollDie(sides)
			if err != nil {
				return nil, err
			}
			out[i] = d
		}
	}

	return out, nil
}
