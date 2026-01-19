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
