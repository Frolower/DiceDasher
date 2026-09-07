package vtmv5

import "github.com/google/uuid"

type rerollRequest struct {
	RecordID    uuid.UUID `json:"record_id"`
	RerollIndex []int     `json:"reroll_index"`
}
type rerollResponse struct {
	state            rollState
	Expression       string `json:"expression"`
	RerollExpression string `json:"reroll_expression"`
	MainRoll         []int  `json:"main_roll"`
	HungerRoll       []int  `json:"hunger_roll"`
	Successes        int    `json:"successes"`
	Success          bool   `json:"success"`
	IsCritical       bool   `json:"is_critical"`
	CritType         string `json:"crit_type"`
}

func (r rerollResponse) HistoryState() any { return r.state }
