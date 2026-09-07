package dice

import "fmt"

// Resource limits, independent of any particular game's rules.
const (
	MaxDice  = 1000
	MaxSides = 1_000_000
)

// Pool is an immutable, validated dice pool. Empty pools are valid for optional
// parts of a roll (gear, hunger); callers enforce nonempty complete rolls.
type Pool struct{ count, sides int }

func NewPool(count, sides int) (Pool, error) {
	if count < 0 || count > MaxDice {
		return Pool{}, fmt.Errorf("dice count must be between 0 and %d", MaxDice)
	}
	if sides < 2 || sides > MaxSides {
		return Pool{}, fmt.Errorf("sides must be between 2 and %d", MaxSides)
	}
	return Pool{count: count, sides: sides}, nil
}
func (p Pool) Count() int           { return p.count }
func (p Pool) Roll() ([]int, error) { return RollDice(p.count, p.sides) }

// AddCounts rejects integer overflow before evaluating each addition.
func AddCounts(values ...int) (int, error) {
	max := int(^uint(0) >> 1)
	min := -max - 1
	total := 0
	for _, value := range values {
		if (value > 0 && total > max-value) || (value < 0 && total < min-value) {
			return 0, fmt.Errorf("dice count overflow")
		}
		total += value
	}
	return total, nil
}
