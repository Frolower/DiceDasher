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

func ResolveTES(r model.ResolveTES) model.ResolveTESResponse {
	res := model.ResolveTESResponse{}

	res.Expression = fmt.Sprintf("%dd6", r.Attr+r.Assist+r.Gear+r.Modificator)
	res.Rolls = dice.RollDice(r.Attr+r.Assist+r.Gear+r.Modificator, 6)
	res.Successes = countInt(res.Rolls, 6)
	if res.Successes < r.Target {
		res.Success = false
	} else {
		res.Success = true
	}

	return res
}

//func ResolveVtMV5r (r model.ResolveVtMV5) model.ResolveVtMV5Response {}

func sum(arr []int) int {
	res := 0
	for _, e := range arr {
		res += e
	}
	return res
}

func countInt(arr []int, target int) int {
	res := 0
	for _, e := range arr {
		if e == target {
			res++
		}
	}
	return res
}
