package tes

import "github.com/google/uuid"

type pushRequest struct {
	RecordID uuid.UUID `json:"record_id"`
}
type pushResponse struct {
	state          rollState
	Expression     string `json:"expression"`
	PushExpression string `json:"push_expression"`
	AttributeRolls []int  `json:"attribute_rolls"`
	GearRolls      []int  `json:"gear_rolls"`
	Successes      int    `json:"successes"`
	Success        bool   `json:"success"`
	HopeLosses     int    `json:"hope_losses"`
	GearDamage     int    `json:"gear_damage"`
}

func (r pushResponse) HistoryState() any { return r.state }
