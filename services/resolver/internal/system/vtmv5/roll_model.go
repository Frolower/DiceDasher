package vtmv5

type rollRequest struct {
	Attribute int `json:"attribute"`
	Skill     int `json:"skill"`
	Hunger    int `json:"hunger"`
	Target    int `json:"target"`
}

type rollResponse struct {
	Expression string `json:"expression"`
	MainRoll   []int  `json:"main_roll"`
	HungerRoll []int  `json:"hunger_roll"`
	Successes  int    `json:"successes"`
	Success    bool   `json:"success"`
	IsCritical bool   `json:"is_critical"`
	CritType   string `json:"crit_type"`
}
