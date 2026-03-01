package model

type Resolve struct {
	Number int `json:"number"`
	Size   int `json:"size"`
}

type ResolveResponse struct {
	Expression string `json:"expression"`
	Rolls      []int  `json:"rolls"`
	Sum        int    `json:"sum"`
}

type ResolveTES struct {
	Attr        int `json:"attr"`        // Main STAT attribute score
	Assist      int `json:"assist"`      // Number of players assisting
	Gear        int `json:"gear"`        // Modificator from the gear
	Modificator int `json:"modificator"` // Extra modificator from GM
	Target      int `json:"target"`      // Target number of successes
}

type ResolveTESResponse struct {
	Expression string `json:"expression"`
	Rolls      []int  `json:"rolls"`
	Successes  int    `json:"successes"`
	Success    bool   `json:"success"`
}

//type ResolveVtMV5 {}

//type ResolveVtMV5Response struct {}
