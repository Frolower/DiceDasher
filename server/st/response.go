package st

import (
	"encoding/json"
)

type Response struct {
	Room_id string `json:"room_id"`
	Player_id string `json:"player_id"`
	Action string `json:"action"`
	Data string `json:"data"`
}


func (g Response) JSON() []byte {
	res, _ := json.Marshal(g)
	// TODO: Handle error 
	return res
}