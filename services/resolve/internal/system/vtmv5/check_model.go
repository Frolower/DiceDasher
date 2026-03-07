package vtmv5

type checkResponse struct {
	Expression string `json:"expression"`
	Result     int    `json:"result"`
	Success    bool   `json:"success"`
}
