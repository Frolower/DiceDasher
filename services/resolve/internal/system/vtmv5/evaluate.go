package vtmv5

import "diceDasher/pkg/util"

// Outcome is calculated from valid d10 results without I/O or randomness.
type Outcome struct {
	Successes  int
	Success    bool
	IsCritical bool
	CritType   string
}

// Evaluate preserves the existing critical classification, including a pair of
// tens taking precedence over bestial failure even if the target is not reached.
// It does not modify either input slice.
func Evaluate(mainRoll, hungerRoll []int, target int) Outcome {
	regular10s := util.CountInt(mainRoll, 10)
	hunger10s := util.CountInt(hungerRoll, 10)
	successes := util.CountBetween(mainRoll, 6, 10) + util.CountBetween(hungerRoll, 6, 10)
	pairsOf10s := (regular10s + hunger10s) / 2
	successes += pairsOf10s * 2 // Adds 2 extra successes for each pair of 10's
	success := successes >= target
	hungerFails := util.CountInt(hungerRoll, 1)
	isCritical := false
	critType := "none"

	if !success && hungerFails > 0 {
		isCritical = true
		critType = "bestial failure"
	}

	if pairsOf10s > 0 {
		isCritical = true
		if hunger10s > 0 {
			critType = "messy critical"
		} else {
			critType = "critical"
		}
	}

	return Outcome{Successes: successes, Success: success, IsCritical: isCritical, CritType: critType}
}
