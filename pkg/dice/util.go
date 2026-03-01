package dice

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
