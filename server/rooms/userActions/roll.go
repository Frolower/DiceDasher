package userActions

import (
	"math/rand"
	"time"
)

func RollDice(size int) []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(size) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(4)+1)
	}
	return result
}