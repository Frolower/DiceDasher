package generic

type request struct {
	Number int `json:"number"`
	Size   int `json:"size"`
}

type response struct {
	Expression string `json:"expression"`
	Rolls      []int  `json:"rolls"`
	Sum        int    `json:"sum"`
}
