package tes

type request struct {
	Attr        int `json:"attr"`        // Main STAT attribute score
	Assist      int `json:"assist"`      // Number of players assisting
	Gear        int `json:"gear"`        // Modificator from the gear
	Modificator int `json:"modificator"` // Extra modificator from GM
	Target      int `json:"target"`      // Target number of successes
}

type response struct {
	Expression string `json:"expression"`
	Rolls      []int  `json:"rolls"`
	Successes  int    `json:"successes"`
	Success    bool   `json:"success"`
}
