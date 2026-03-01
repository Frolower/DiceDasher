package dice

import (
	"fmt"
	"math/rand/v2"
)

func RollDie(sides int) (int, error) {
	if sides < 2 {
		return 0, fmt.Errorf("sides must be >= 2")
	}
	roll := rand.IntN(sides) + 1
	return roll, nil
}

func RollDice(count int, sides int) ([]int, error) {
	if sides < 2 {
		return nil, fmt.Errorf("sides must be >= 2")
	}

	rolls := make([]int, count)
	for i := 0; i < count; i++ {
		rolls[i] = rand.IntN(sides) + 1
	}
	return rolls, nil
}

func RerollKeepingValues(previous []int, except []int, sides int) ([]int, error) {
	if sides < 2 {
		return nil, fmt.Errorf("sides must be >= 2")
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
