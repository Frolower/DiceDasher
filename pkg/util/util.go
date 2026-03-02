package util

func Sum(arr []int) int {
	res := 0
	for _, e := range arr {
		res += e
	}
	return res
}

func CountInt(arr []int, target int) int {
	res := 0
	for _, e := range arr {
		if e == target {
			res++
		}
	}
	return res
}

func CountBetween(arr []int, low, high int) int {
	res := 0
	for _, e := range arr {
		if e >= low && e <= high {
			res++
		}
	}
	return res
}

func CountAbove(arr []int, target int) int {
	res := 0
	for _, e := range arr {
		if e > target {
			res++
		}
	}
	return res
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
