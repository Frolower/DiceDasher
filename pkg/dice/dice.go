package dice

import "math/rand/v2"

func RollDice(count int, sides int) []int {
	rolls := make([]int, count)
	for i := 0; i < count; i++ {
		rolls[i] = rand.IntN(sides) + 1
	}
	return rolls
}
