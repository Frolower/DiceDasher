package tes

import "github.com/google/uuid"

type pushRequest struct {
	RecordID uuid.UUID `json:"record_id"`
}
type pushRecord struct {
	AttributeRolls []int `json:"attribute_rolls"`
	GearRolls      []int `json:"gear_rolls"`
	Target         int   `json:"target"`
}

type pushResponse struct {
	Expression     string `json:"expression"`
	PushExpression string `json:"push_expression"`
	AttributeRolls []int  `json:"attribute_rolls"`
	GearRolls      []int  `json:"gear_rolls"`
	Successes      int    `json:"successes"`
	Success        bool   `json:"success"`
	HopeLosses     int    `json:"hope_losses"`
	GearDamage     int    `json:"gear_damage"`
}
