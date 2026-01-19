package service

import (
	"diceDasher/pkg/dice"
	"diceDasher/services/resolve/internal/model"
	"fmt"
)

func ResolveRoll(r model.Resolve) model.ResolveResponse {
	res := model.ResolveResponse{}

	res.Expression = fmt.Sprintf("%dd%d", r.Number, r.Size)
	res.Rolls = dice.RollDice(r.Number, r.Size)
	res.Sum = sum(res.Rolls)
	return res
}

func sum(arr []int) int {
	res := 0
	for _, e := range arr {
		res += e
	}
	return res
}
