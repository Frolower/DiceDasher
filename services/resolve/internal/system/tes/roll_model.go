package tes

type rollRequest struct {
	Attr        int `json:"attr"`        // Main STAT attribute score
	Assist      int `json:"assist"`      // Number of players assisting
	Gear        int `json:"gear"`        // Modificator from the gear
	Modificator int `json:"modificator"` // Extra modificator from GM
	Target      int `json:"target"`      // Target number of successes
}

type rollResponse struct {
	Expression     string `json:"expression"`
	AttributeRolls []int  `json:"attribute_rolls"`
	GearRolls      []int  `json:"gear_rolls"`
	Successes      int    `json:"successes"`
	Success        bool   `json:"success"`
}
