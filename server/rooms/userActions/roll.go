package userActions

import (
	"math/rand"
	"time"
)

func rollD2() []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(6) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(2)+1)
	}
	return result
}

func rollD4() []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(6) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(4)+1)
	}
	return result
}

func rollD6() []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(6) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(6)+1)
	}
	return result
}

func rollD8() []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(6) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(8)+1)
	}
	return result
}

func rollD10() []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(6) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(10)+1)
	}
	return result
}

func rollD12() []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(6) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(12)+1)
	}
	return result
}

func rollD20() []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(6) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(20)+1)
	}
	return result
}

func rollD100() []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(6) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(100)+1)
	}
	return result
}
