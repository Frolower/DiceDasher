package rooms

import (
	"encoding/json"
)

type Response struct {
	Room_id string `json:"room_id"`
	Player_id string `json:"player_id"`
	Action string `json:"action"`
	Data map [string]string `json:"data"`
}


func (this Response) JSON() []byte {
	res, _ := json.Marshal(this)
	// TODO: Handle error 
	return res
}