package model

type DiceRollRequest struct {
	Type   int `json:"type"`
	Amount int `json:"amount"`
}

type DiceRollResult struct {
	Type    int    `json:"type"`
	Results []int  `json:"results"`
	IsCrit  []bool `json:"is_crit"`
	Sum     int    `json:"sum"`
}

type RollAgainstRequest struct {
	Type        int `json:"type"`
	Amount      int `json:"amount"`
	RollAgainst int `json:"roll_against"`
}

type RollAgainstResult struct {
	Type           int   `json:"type"`
	Results        []int `json:"results"`
	ResultsAgainst []bool
	IsCrit         []bool `json:"is_crit"`
	Sum            int    `json:"sum"`
}
